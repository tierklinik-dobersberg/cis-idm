package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func (svc *Service) AddPhoneNumber(ctx context.Context, req *connect.Request[idmv1.AddPhoneNumberRequest]) (*connect.Response[idmv1.AddPhoneNumberResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	m := models.PhoneNumber{
		UserID:      claims.Subject,
		PhoneNumber: req.Msg.Number,
		Verified:    false,
		Primary:     false,
	}

	m, err := svc.Datastore.AddUserPhoneNumber(ctx, m)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.AddPhoneNumberResponse{
		PhoneNumber: conv.PhoneNumberProtoFromPhoneNumber(m),
	}), nil
}

func (svc *Service) DeletePhoneNumber(ctx context.Context, req *connect.Request[idmv1.DeletePhoneNumberRequest]) (*connect.Response[idmv1.DeletePhoneNumberResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Datastore.DeleteUserPhoneNumber(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.DeletePhoneNumberResponse{}), nil
}

func (svc *Service) MarkPhoneNumberAsPrimary(ctx context.Context, req *connect.Request[idmv1.MarkPhoneNumberAsPrimaryRequest]) (*connect.Response[idmv1.MarkPhoneNumberAsPrimaryResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Datastore.MarkPhoneNumberAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkPhoneNumberAsPrimaryResponse{}), nil
}
