package common

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	"github.com/gofrs/uuid"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

var (
	ErrFeatureDisabled = connect.NewError(connect.CodeUnavailable, fmt.Errorf("the requested feature has been disabled by an administrator"))
)

func (svc *Service) AddUserAddress(ctx context.Context, model repo.UserAddress) ([]repo.UserAddress, error) {
	if svc.cfg.DisableUserAddresses {
		return nil, ErrFeatureDisabled
	}

	if model.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		model.ID = id.String()
	}

	if _, err := svc.repo.CreateUserAddress(ctx, repo.CreateUserAddressParams{
		ID:       model.ID,
		UserID:   model.UserID,
		CityCode: model.CityCode,
		CityName: model.CityName,
		Street:   model.Street,
		Extra:    model.Extra,
	}); err != nil {
		return nil, fmt.Errorf("failed to save new user address: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, model.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
	}

	return addresses, nil
}

func (svc *Service) DeleteUserAddress(ctx context.Context, userID string, addressID string) ([]repo.UserAddress, error) {
	if svc.cfg.DisableUserAddresses {
		return nil, ErrFeatureDisabled
	}

	if rows, err := svc.repo.DeleteUserAddress(ctx, repo.DeleteUserAddressParams{
		ID:     addressID,
		UserID: userID,
	}); err == nil {
		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("address not fund"))
		}
	} else {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
	}

	return addresses, nil
}

func (svc *Service) UpdateUserAddress(ctx context.Context, updateModel repo.UserAddress, paths []string) ([]repo.UserAddress, error) {
	if svc.cfg.DisableUserAddresses {
		return nil, ErrFeatureDisabled
	}

	addr, err := svc.repo.GetUserAddress(ctx, repo.GetUserAddressParams{
		UserID: updateModel.UserID,
		ID:     updateModel.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("address does not exist"))
		}

		return nil, fmt.Errorf("failed to load address by id: %w", err)
	}

	if len(paths) == 0 {
		paths = []string{
			"city_code", "city", "street", "extra",
		}
	}

	// create a "add-address-request" that we can use to validate the update
	// operation.
	protoAddr := &idmv1.AddAddressRequest{
		CityCode: addr.CityCode,
		CityName: addr.CityName,
		Extra:    addr.Extra,
		Street:   addr.Street,
	}

	for _, p := range paths {
		switch p {
		case "city_code":
			addr.CityCode = updateModel.CityCode
			protoAddr.CityCode = updateModel.CityCode
		case "city_name":
			addr.CityName = updateModel.CityName
			protoAddr.CityName = updateModel.CityName
		case "street":
			addr.Street = updateModel.Street
			protoAddr.Street = updateModel.Street
		case "extra":
			addr.Extra = updateModel.Extra
			protoAddr.Extra = updateModel.Extra
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid field_mask for update operation: invalid path %q", p))
		}
	}

	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	if err := validator.Validate(protoAddr); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if _, err := svc.repo.UpdateUserAddress(ctx, repo.UpdateUserAddressParams{
		CityCode: addr.CityCode,
		CityName: addr.CityName,
		Street:   addr.Street,
		Extra:    addr.Extra,
		ID:       addr.ID,
		UserID:   addr.UserID,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user address: %w", err))
	}

	addrs, err := svc.repo.GetUserAddresses(ctx, updateModel.UserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load user addresses: %w", err))
	}

	return addrs, nil
}
