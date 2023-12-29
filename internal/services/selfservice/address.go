package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func (svc *Service) AddAddress(ctx context.Context, req *connect.Request[idmv1.AddAddressRequest]) (*connect.Response[idmv1.AddAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	addresses, err := svc.Common.AddUserAddress(ctx, repo.UserAddress{
		UserID:   claims.Subject,
		CityCode: req.Msg.CityCode,
		CityName: req.Msg.CityName,
		Street:   req.Msg.Street,
		Extra:    req.Msg.Extra,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save new user address: %w", err)
	}

	return connect.NewResponse(&idmv1.AddAddressResponse{
		Addresses: conv.AddressProtosFromAddresses(addresses...),
	}), nil
}

func (svc *Service) DeleteAddress(ctx context.Context, req *connect.Request[idmv1.DeleteAddressRequest]) (*connect.Response[idmv1.DeleteAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	addresses, err := svc.Common.DeleteUserAddress(ctx, claims.Subject, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.DeleteAddressResponse{
		Addresses: conv.AddressProtosFromAddresses(addresses...),
	}), nil
}

func (svc *Service) UpdateAddress(ctx context.Context, req *connect.Request[idmv1.UpdateAddressRequest]) (*connect.Response[idmv1.UpdateAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	addrs, err := svc.Common.UpdateUserAddress(ctx, repo.UserAddress{
		CityCode: req.Msg.CityCode,
		CityName: req.Msg.CityName,
		Street:   req.Msg.Street,
		Extra:    req.Msg.Extra,
		UserID:   claims.Subject,
		ID:       req.Msg.Id,
	}, req.Msg.FieldMask.Paths)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.UpdateAddressResponse{
		Addresses: conv.AddressProtosFromAddresses(addrs...),
	}), nil
}
