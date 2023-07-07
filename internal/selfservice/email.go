package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func (svc *Service) AddEmailAddress(ctx context.Context, req *connect.Request[idmv1.AddEmailAddressRequest]) (*connect.Response[idmv1.AddEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	mails, err := svc.Common.AddEmailAddressToUser(ctx, models.EMail{
		UserID:  claims.Subject,
		Address: req.Msg.Email,
	})
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&idmv1.AddEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) DeleteEmailAddress(ctx context.Context, req *connect.Request[idmv1.DeleteEmailAddressRequest]) (*connect.Response[idmv1.DeleteEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	mails, err := svc.Common.DeleteEmailAddressFromUser(ctx, claims.Subject, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&idmv1.DeleteEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) MarkEmailAsPrimary(ctx context.Context, req *connect.Request[idmv1.MarkEmailAsPrimaryRequest]) (*connect.Response[idmv1.MarkEmailAsPrimaryResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Common.MarkEmailAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkEmailAsPrimaryResponse{}), nil
}
