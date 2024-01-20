package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/SherClockHolmes/webpush-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

func (svc *Service) sendWebPushNotification(
	ctx context.Context,
	user *idmv1.Profile,
	tmplCtx map[string]any,
	msg *idmv1.SendNotificationRequest_Webpush) ([]*idmv1.DeliveryNotification, error) {

	userID := user.User.Id

	subscriptions, err := svc.Datastore.GetWebPushSubscriptionsForUser(ctx, userID)
	if err != nil {
		return []*idmv1.DeliveryNotification{
			{
				TargetUser: user.User.Id,
				Error:      err.Error(),
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_OTHER,
			},
		}, nil
	}

	if len(subscriptions) == 0 {
		return []*idmv1.DeliveryNotification{
			{
				TargetUser: user.User.Id,
				Error:      "user does not have web-push enabled",
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_NO_WEBPUSH_SUBSCRIPTION,
			},
		}, nil
	}

	var (
		deliveries []*idmv1.DeliveryNotification
		content    []byte
	)

	switch kind := msg.Webpush.Kind.(type) {
	case *idmv1.WebPushNotification_Binary:
		content = kind.Binary

	case *idmv1.WebPushNotification_Notification:

		notif := notifFromProto(kind.Notification)

		if tmplCtx != nil {
			t := textTemplate.New("")

			fm := textTemplate.FuncMap(tmpl.PrepareFunctionMap())
			tmpl.AddToMap(fm, sprig.TxtFuncMap())

			t.Funcs(fm)

			if _, err := t.Parse(notif.Body); err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)

			if err := t.Execute(buf, tmplCtx); err == nil {

				notif.Body = buf.String()

			} else {
				deliveries = append(deliveries, &idmv1.DeliveryNotification{
					TargetUser: userID,
					Error:      fmt.Errorf("failed to execute template: %w", err).Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TEMPLATE,
				})

				return deliveries, nil
			}
		}

		blob, err := json.Marshal(map[string]any{
			"notification": notif,
		})
		if err != nil {
			deliveries = append(deliveries, &idmv1.DeliveryNotification{
				TargetUser: userID,
				Error:      fmt.Errorf("failed to execute template: %w", err).Error(),
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TEMPLATE,
			})

			return deliveries, nil
		}

		content = blob

	default:
		deliveries = append(deliveries, &idmv1.DeliveryNotification{
			TargetUser: userID,
			Error:      fmt.Errorf("unsupported payload: %T", kind).Error(),
			ErrorKind:  idmv1.ErrorKind_ERROR_KIND_OTHER,
		})

		return deliveries, nil
	}

	log.L(ctx).Infof("sending web-push notification to %s (username=%q)", user.User.Id, user.User.Username)

	atLeastOneSuccess := false
	for _, sub := range subscriptions {
		webpushSub := webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				Auth:   sub.Auth,
				P256dh: sub.Key,
			},
		}

		res, err := webpush.SendNotificationWithContext(ctx, content, &webpushSub, &webpush.Options{
			Subscriber:      svc.Config.WebPush.Admin,
			VAPIDPublicKey:  svc.Config.WebPush.VAPIDpublicKey,
			VAPIDPrivateKey: svc.Config.WebPush.VAPIDprivateKey,
		})
		if err != nil {
			deliveries = append(deliveries, &idmv1.DeliveryNotification{
				TargetUser: userID,
				Error:      fmt.Errorf("transport error: %w", err).Error(),
				ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TRANSPORT,
			})

			continue
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			if res.StatusCode == http.StatusGone {
				// this subscription has been removed/expired so remove it from the db
				go func(sub repo.WebpushSubscription) {
					rows, err := svc.Datastore.DeleteWebPushSubscriptionByID(context.Background(), sub.ID)
					if err != nil {
						log.L(context.Background()).Errorf("failed to delete expired web-push subscription %s: %s", sub.ID, err)
					} else if rows == 0 {
						log.L(context.Background()).Errorf("failed to delete expired web-push subscription %s: not found", sub.ID)
					}
				}(sub)
			} else {
				deliveries = append(deliveries, &idmv1.DeliveryNotification{
					TargetUser: userID,
					Error:      fmt.Errorf("unexpected response from web-push service: status-code=%d", res.StatusCode).Error(),
					ErrorKind:  idmv1.ErrorKind_ERROR_KIND_TRANSPORT,
				})
			}

			continue
		} else {
			atLeastOneSuccess = true
		}
	}

	if atLeastOneSuccess {
		// successfull delivery
		deliveries = append(deliveries, &idmv1.DeliveryNotification{
			TargetUser: userID,
		})
	}

	return deliveries, nil
}

type action struct {
	Action string `json:"action"`
	Title  string `json:"title"`
}

type swNotification struct {
	Title   string         `json:"title,omitempty"`
	Body    string         `json:"body,omitempty"`
	Actions []action       `json:"actions,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type operation struct {
	Operation string `json:"operation,omitempty"`
	URL       string `json:"url,omitempty"`
}

func notifFromProto(pb *idmv1.ServiceWorkerNotification) swNotification {
	n := swNotification{
		Title: pb.Title,
		Body:  pb.Body,
	}

	handlers := make(map[string]operation)

	for _, pbAction := range pb.Actions {
		handlers[pbAction.Action] = operation{
			Operation: operationFromProto(pbAction.Operation),
			URL:       pbAction.Url,
		}

		n.Actions = append(n.Actions, action{
			Action: pbAction.Action,
			Title:  pbAction.Title,
		})
	}

	if pb.DefaultOperation != idmv1.Operation_OPERATION_UNSPECIFIED {
		handlers["default"] = operation{
			Operation: operationFromProto(pb.DefaultOperation),
			URL:       pb.DefaultOperationUrl,
		}
	}

	n.Data = map[string]any{
		"onActionClick": handlers,
	}

	return n
}

func operationFromProto(op idmv1.Operation) string {
	switch op {
	case idmv1.Operation_OPERATION_FOCUS_LAST_FOCUSED_OR_OPEN:
		return "focusLastFocusedOrOpen"
	case idmv1.Operation_OPERATION_NAVIGATE_LAST_FOCUSED_OR_OPEN:
		return "navigateLastFocusedOrOpen"
	case idmv1.Operation_OPERATION_OPEN_WINDOW:
		return "openWindow"
	case idmv1.Operation_OPERATION_SEND_REQUEST:
		return "sendRequest"
	}

	return ""
}
