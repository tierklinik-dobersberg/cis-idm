package notify

import (
	"bytes"
	"context"
	"fmt"
	htmlTemplate "html/template"
	"io"

	"github.com/Masterminds/sprig/v3"
	"github.com/ory/mail"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"golang.org/x/exp/slices"
)

func (svc *Service) sendMail(
	ctx context.Context,
	user *idmv1.Profile,
	tmplCtx map[string]any,
	m *mail.Message,
	msg *idmv1.SendNotificationRequest_Email) ([]*idmv1.DeliveryNotification, error) {

	if svc.Config.MailConfig == nil || svc.Config.MailConfig.Host == "" {
		return nil, fmt.Errorf("SMTP not configured")
	}

	userID := user.User.Id

	addr := user.User.GetPrimaryMail().GetAddress()
	if addr == "" {
		return []*idmv1.DeliveryNotification{
			{
				TargetUser: user.User.Id,
				Error:      "user does not have a primary mail",
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_NO_PRIMARY_MAIL,
			},
		}, nil
	}

	log.L(ctx).Info("preparing mail", "userId", user.User.Id, "username", user.User.Username, "address", addr)

	if dr := svc.Config.DryRun; dr != nil && dr.MailTarget != "" {
		log.L(ctx).Info("replacing receipient address dry-run mode", "old", addr, "new", dr.MailTarget)

		addr = dr.MailTarget
	}

	m.SetHeaders(map[string][]string{
		"From":    {svc.Config.MailConfig.From},
		"To":      {addr},
		"Subject": {msg.Email.Subject},
	})

	buf := new(bytes.Buffer)
	if tmplCtx == nil {

		_, err := buf.Write([]byte(msg.Email.Body))
		if err != nil {
			return []*idmv1.DeliveryNotification{
				{
					TargetUser: userID,
					Error:      fmt.Sprintf("failed to write email body to buffer: %s", err),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_OTHER,
				},
			}, nil
		}

	} else {
		t := htmlTemplate.New("")

		fm := htmlTemplate.FuncMap(tmpl.PrepareFunctionMap(svc.Datastore))
		tmpl.AddToMap(fm, sprig.HtmlFuncMap())

		t.Funcs(fm)

		if _, err := t.Parse(msg.Email.Body); err != nil {
			return []*idmv1.DeliveryNotification{
				{
					TargetUser: userID,
					Error:      err.Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TEMPLATE,
				},
			}, nil
		}

		if err := t.Execute(buf, tmplCtx); err != nil {
			return []*idmv1.DeliveryNotification{
				{
					TargetUser: userID,
					Error:      err.Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TEMPLATE,
				},
			}, nil
		}
	}

	m.SetBodyWriter("text/html", func(w io.Writer) error {
		_, err := w.Write(buf.Bytes())
		return err
	})

	for _, at := range msg.Email.Attachments {
		if len(at.ForUser) > 0 {
			if !slices.Contains(at.ForUser, userID) {
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
			return []*idmv1.DeliveryNotification{
				{
					TargetUser: userID,
					Error:      "unsupported attachment type",
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_OTHER,
				},
			}, nil
		}
	}

	log.L(ctx).Info("sending mails")
	if err := svc.Mailer.DialAndSend(m); err != nil {
		return []*idmv1.DeliveryNotification{
			{
				TargetUser: userID,
				Error:      err.Error(),
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TRANSPORT,
			},
		}, nil
	}

	// Successfully sent the mail
	return []*idmv1.DeliveryNotification{{TargetUser: userID}}, nil
}
