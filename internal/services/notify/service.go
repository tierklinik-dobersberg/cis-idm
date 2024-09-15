package notify

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	"github.com/ory/mail"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
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

func (svc *Service) getSenderUser(ctx context.Context, id string) (*idmv1.Profile, error) {
	if id == "" {
		claims := middleware.ClaimsFromContext(ctx)
		if claims == nil {
			return nil, fmt.Errorf("no claims associated with context")
		}

		id = claims.Subject
	}

	log.L(ctx).Debugf("loading sender user model for user %q", id)
	senderUserModel, err := svc.Datastore.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load user object for sender %q: %w", id, err)
	}

	senderUser, err := svc.GetUserProfileProto(ctx, senderUserModel)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile of sender %q: %w", senderUserModel.ID, err)
	}

	return senderUser, nil
}

func (svc *Service) loadUsers(ctx context.Context, userIds []string) (map[string]*idmv1.Profile, error) {
	targetUsers := make(map[string]*idmv1.Profile)
	for _, usr := range userIds {
		userModel, err := svc.Datastore.GetUserByID(ctx, usr)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user with id %q not found", usr))
			}

			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load user with id %q: %w", usr, err))
		}

		profile, err := svc.GetUserProfileProto(ctx, userModel)
		if err != nil {
			return nil, fmt.Errorf("failed to load user: %s (username=%q): %w", userModel.ID, userModel.Username, err)
		}

		targetUsers[usr] = profile
	}

	return targetUsers, nil
}

func (svc *Service) SendNotification(ctx context.Context, req *connect.Request[idmv1.SendNotificationRequest]) (*connect.Response[idmv1.SendNotificationResponse], error) {
	log.L(ctx).Infof("received SendNotification request")

	var senderUser *idmv1.Profile
	if req.Msg.SenderUserId != "" {
		var err error
		senderUser, err = svc.getSenderUser(ctx, req.Msg.SenderUserId)
		if err != nil {
			return nil, err
		}
	}

	targetUsers, err := svc.loadUsers(ctx, req.Msg.TargetUsers)
	if err != nil {
		return nil, err
	}

	if len(targetUsers) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("no receipients specified"))
	}

	log.L(ctx).Infof("gathering user information for %d receipients", len(targetUsers))

	prepareTmplContext := func(usrID string, ctxPb *structpb.Struct, msg *mail.Message) map[string]any {
		var m map[string]any
		if ctxPb == nil {
			m = make(map[string]any)
		} else {
			m = ctxPb.AsMap()
		}

		rctx := tmpl.NewRenderContext()
		rctx.Set("mail", msg)
		rctx.Set("ctx", ctx)

		m["User"] = targetUsers[usrID]
		m["Sender"] = senderUser
		m["IDM"] = map[string]any{
			"SiteName":  svc.Config.UserInterface.SiteName,
			"SiteURL":   svc.Config.UserInterface.SiteNameURL,
			"PublicURL": svc.Config.UserInterface.PublicURL,
		}
		m["Ctx"] = rctx

		return m
	}

	var deliveries []*idmv1.DeliveryNotification

	switch msg := req.Msg.Message.(type) {
	case *idmv1.SendNotificationRequest_Email:
		for userID, profile := range targetUsers {

			m := mail.NewMessage()
			tmplCtx := prepareTmplContext(userID, req.Msg.PerUserTemplateContext[userID], m)

			result, err := svc.sendMail(ctx, profile, tmplCtx, m, msg)
			if err != nil {
				return nil, err
			}

			deliveries = append(deliveries, result...)
		}

	case *idmv1.SendNotificationRequest_Sms:
		for userID, profile := range targetUsers {
			addr := profile.User.GetPrimaryPhoneNumber().GetNumber()
			if addr == "" {
				deliveries = append(deliveries, &idmv1.DeliveryNotification{
					TargetUser: userID,
					Error:      "user does not have a primary phone number",
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_NO_PRIMARY_PHONE,
				})

				continue
			}

			body := msg.Sms.Body

			tmplCtx := prepareTmplContext(userID, req.Msg.PerUserTemplateContext[userID], nil)

			var err error
			if tmplCtx != nil {
				log.L(ctx).Infof("preparing template for user %s to phone number %s", userID, addr)

				t := textTemplate.New("")

				fm := textTemplate.FuncMap(tmpl.PrepareFunctionMap(svc.Datastore))
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

			if err != nil {
				deliveries = append(deliveries, &idmv1.DeliveryNotification{
					TargetUser: userID,
					Error:      err.Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TEMPLATE,
				})

				continue
			}

			log.L(ctx).Infof("sending SMS to %s", addr)
			if err = svc.SMSSender.Send(ctx, sms.Message{
				From: svc.Config.Twilio.From,
				To:   []string{addr},
				Body: body,
			}); err != nil {
				deliveries = append(deliveries, &idmv1.DeliveryNotification{
					TargetUser: userID,
					Error:      err.Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TRANSPORT,
				})

				continue
			}

			deliveries = append(deliveries, &idmv1.DeliveryNotification{
				TargetUser: userID,
			})
		}

	case *idmv1.SendNotificationRequest_Webpush:
		if svc.Config.WebPush == nil {
			return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("web-push not configured"))
		}

		for userID, profile := range targetUsers {
			tmplCtx := prepareTmplContext(userID, req.Msg.PerUserTemplateContext[userID], nil)

			result, err := svc.sendWebPushNotification(ctx, profile, tmplCtx, msg)
			if err != nil {
				return nil, err
			}

			deliveries = append(deliveries, result...)
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing notification type"))
	}

	return connect.NewResponse(&idmv1.SendNotificationResponse{
		Deliveries: deliveries,
	}), nil
}

func (svc *Service) GetVAPIDPublicKey(ctx context.Context, req *connect.Request[idmv1.GetVAPIDPublicKeyRequest]) (*connect.Response[idmv1.GetVAPIDPublicKeyResponse], error) {
	if svc.Config.WebPush == nil {
		return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("web-push unconfigured"))
	}

	return connect.NewResponse(&idmv1.GetVAPIDPublicKeyResponse{
		Key: svc.Config.WebPush.VAPIDpublicKey,
	}), nil
}

func (svc *Service) AddWebPushSubscription(ctx context.Context, req *connect.Request[idmv1.AddWebPushSubscriptionRequest]) (*connect.Response[idmv1.AddWebPushSubscriptionResponse], error) {
	if svc.Config.WebPush == nil {
		return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("web-push unconfigured"))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	tokenID := claims.ID
	if claims.AppMetadata != nil {
		tokenID = claims.AppMetadata.ParentTokenID
	}

	keys := webpush.Keys{
		Auth:   req.Msg.Subscription.Keys.Auth,
		P256dh: req.Msg.Subscription.Keys.P256Dh,
	}

	sub := webpush.Subscription{
		Endpoint: req.Msg.Subscription.Endpoint,
		Keys:     keys,
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	params := repo.CreateWebPushSubscriptionForUserParams{
		ID:        id.String(),
		UserID:    claims.Subject,
		UserAgent: req.Header().Get("User-Agent"),
		Endpoint:  sub.Endpoint,
		Auth:      sub.Keys.Auth,
		Key:       sub.Keys.P256dh,
		TokenID:   tokenID,
	}

	if err := svc.Datastore.CreateWebPushSubscriptionForUser(ctx, params); err != nil {
		return nil, err
	}

	return connect.NewResponse(new(idmv1.AddWebPushSubscriptionResponse)), nil
}
