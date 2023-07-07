package notify

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
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
	// TODO
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("SendNotification is not yet implemented"))
}
