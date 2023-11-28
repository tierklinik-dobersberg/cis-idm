package notify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	htmlTemplate "html/template"
	"io"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/bufbuild/connect-go"
	"github.com/ory/mail"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	idmv1connect.UnimplementedNotifyServiceHandler

	*app.Providers
}

func New(providers *app.Providers) *Service {
	return &Service{
		Providers: providers,
	}
}

func (svc *Service) SendNotification(ctx context.Context, req *connect.Request[idmv1.SendNotificationRequest]) (*connect.Response[idmv1.SendNotificationResponse], error) {
	log.L(ctx).Infof("received SendNotification request")
	targetUsers := make(map[string]*idmv1.Profile)

	senderUserModel, err := svc.Datastore.GetUserByID(ctx, req.Msg.SenderUserId)
	if err != nil {
		return nil, fmt.Errorf("failed to load user object for sender %q: %w", req.Msg.SenderUserId, err)
	}

	senderUser, err := svc.GetUserProfileProto(ctx, senderUserModel)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile of sender %q: %w", senderUserModel.ID, err)
	}

	for _, usr := range req.Msg.TargetUsers {
		userModel, err := svc.Datastore.GetUserByID(ctx, usr)
		if err != nil {
			if errors.Is(err, stmts.ErrNoResults) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user with id %q not found", usr))
			}

			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load user with id %q: %w", usr, err))
		}

		profile, err := svc.GetUserProfileProto(ctx, userModel)
		if err != nil {
			return nil, err
		}

		targetUsers[usr] = profile
	}

	if len(targetUsers) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no receipients specified"))
	}

	log.L(ctx).Infof("gathering user information for %d receipients", len(targetUsers))
	targetAddrs := make(map[string]string, len(targetUsers))

	for _, usr := range targetUsers {
		switch req.Msg.Message.(type) {
		case *idmv1.SendNotificationRequest_Email:
			if pm := usr.User.PrimaryMail; pm != nil {
				targetAddrs[usr.User.Id] = pm.Address
			}

		case *idmv1.SendNotificationRequest_Sms:
			if pp := usr.User.PrimaryPhoneNumber; pp != nil {
				targetAddrs[usr.User.Id] = pp.Number
			}
		}
	}

	if len(targetAddrs) == 0 {
		return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("failed to gather receipient addresses"))
	}

	deliveries := make([]*idmv1.DeliveryNotification, 0, len(targetAddrs))

	prepareTmplContext := func(usrID string, ctxPb *structpb.Struct, msg *mail.Message) map[string]any {
		if ctxPb == nil {
			return nil
		}

		m := ctxPb.AsMap()

		rctx := tmpl.NewRenderContext()
		rctx.Set("mail", msg)
		rctx.Set("ctx", ctx)

		m["User"] = targetUsers[usrID]
		m["Sender"] = senderUser
		m["IDM"] = map[string]any{
			"SiteName":  svc.Config.SiteName,
			"SiteURL":   svc.Config.SiteNameURL,
			"PublicURL": svc.Config.PublicURL,
		}
		m["Ctx"] = rctx

		return m
	}

	switch msg := req.Msg.Message.(type) {
	case *idmv1.SendNotificationRequest_Email:
		mails := make([]*mail.Message, 0, len(targetAddrs))

		for uisr, addr := range targetAddrs {
			log.L(ctx).Infof("preparing mail for user %s to address %s", uisr, addr)

			if dr := svc.Config.DryRun; dr != nil && dr.MailTarget != "" {
				log.L(ctx).Infof("replacing receipient address %s with %s in dry-run mode", addr, dr.MailTarget)

				addr = dr.MailTarget
			}

			m := mail.NewMessage()
			m.SetHeaders(map[string][]string{
				"From":    {svc.Config.MailConfig.From},
				"To":      {addr},
				"Subject": {msg.Email.Subject},
			})

			tmplCtx := prepareTmplContext(uisr, req.Msg.PerUserTemplateContext[uisr], m)

			buf := new(bytes.Buffer)
			if tmplCtx == nil {
				_, err := buf.Write([]byte(msg.Email.Body))
				if err != nil {
					return nil, fmt.Errorf("failed to write email body to buffer: %w", err)
				}
			} else {
				t := htmlTemplate.New("")

				fm := htmlTemplate.FuncMap(tmpl.PrepareFunctionMap())
				tmpl.AddToMap(fm, sprig.HtmlFuncMap())

				t.Funcs(fm)

				if _, err := t.Parse(msg.Email.Body); err != nil {
					return nil, fmt.Errorf("failed to parse template: %w", err)
				}

				if err := t.Execute(buf, tmplCtx); err != nil {
					return nil, fmt.Errorf("failed to execute template: %w", err)
				}
			}

			m.SetBodyWriter("text/html", func(w io.Writer) error {
				_, err := w.Write(buf.Bytes())
				return err
			})

			for _, at := range msg.Email.Attachments {
				if len(at.ForUser) > 0 {
					if !slices.Contains(at.ForUser, uisr) {
						continue
					}
				}

				options := []mail.FileSetting{
					mail.SetHeader(map[string][]string{
						"Content-Type": {at.MediaType},
					}),
				}

				if at.ContentId != "" {
					options = append(options, mail.SetHeader(map[string][]string{
						"Content-ID": {fmt.Sprintf("<%s>", at.ContentId)},
					}))
				}

				switch at.AttachmentType {
				case idmv1.AttachmentType_INLINE:
					m.EmbedReader(at.Name, bytes.NewReader(at.Content), options...)

				case idmv1.AttachmentType_ALTERNATIVE_BODY:
					m.AddAlternativeWriter(at.MediaType, func(w io.Writer) error {
						// FIXME(ppacher): add support for template execution here
						// based on the specified contentType.
						_, err := w.Write(at.Content)
						return err
					})

				case idmv1.AttachmentType_ATTACHEMNT:
					m.AttachReader(at.Name, bytes.NewReader(at.Content), options...)

				default:
					return nil, fmt.Errorf("invalid or unsupported attachemnt type: %s", at.AttachmentType.String())
				}
			}

			mails = append(mails, m)
		}

		log.L(ctx).Infof("sending mails")
		if err := svc.Mailer.DialAndSend(mails...); err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to send one or more emails: %w", err))
		}

	case *idmv1.SendNotificationRequest_Sms:
		for usr, addr := range targetAddrs {
			body := msg.Sms.Body

			tmplCtx := prepareTmplContext(usr, req.Msg.PerUserTemplateContext[usr], nil)

			var err error
			if tmplCtx != nil {
				log.L(ctx).Infof("preparing template for user %s to phone number %s", usr, addr)

				t := textTemplate.New("")

				fm := textTemplate.FuncMap(tmpl.PrepareFunctionMap())
				tmpl.AddToMap(fm, sprig.TxtFuncMap())

				t.Funcs(fm)

				if _, err := t.Parse(msg.Sms.Body); err != nil {
					return nil, err
				}

				buf := new(bytes.Buffer)
				err = t.Execute(buf, tmplCtx)

				if err == nil {
					body = buf.String()
				}
			}

			if err == nil {
				log.L(ctx).Infof("sending SMS to %s", addr)
				err = svc.SMSSender.Send(ctx, sms.Message{
					From: svc.Config.Twilio.From,
					To:   []string{addr},
					Body: body,
				})
			}

			errStr := ""
			if err != nil {
				errStr = err.Error()
			}

			deliveries = append(deliveries, &idmv1.DeliveryNotification{
				TargetUser: usr,
				Error:      errStr,
			})
		}
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing notification type"))
	}

	return connect.NewResponse(&idmv1.SendNotificationResponse{
		Deliveries: deliveries,
	}), nil
}
