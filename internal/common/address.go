package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (svc *Service) AddUserAddress(ctx context.Context, model models.Address) ([]models.Address, error) {
	if !svc.cfg.FeatureEnabled(config.FeatureAddresses) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("addresses: %w", config.ErrFeatureDisabled))
	}

	if _, err := svc.repo.AddUserAddress(ctx, model); err != nil {
		return nil, fmt.Errorf("failed to save new user address: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, model.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
	}

	return addresses, nil
}

func (svc *Service) DeleteUserAddress(ctx context.Context, userID string, addressID string) ([]models.Address, error) {
	if !svc.cfg.FeatureEnabled(config.FeatureAddresses) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("addresses: %w", config.ErrFeatureDisabled))
	}

	if err := svc.repo.DeleteUserAddress(ctx, userID, addressID); err != nil {
		if errors.Is(err, stmts.ErrNoRowsAffected) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("address not fund"))
		}
		return nil, fmt.Errorf("failed to delete user address: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
	}

	return addresses, nil
}

func (svc *Service) UpdateUserAddress(ctx context.Context, updateModel models.Address, paths []string) ([]models.Address, error) {
	if !svc.cfg.FeatureEnabled(config.FeatureAddresses) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("addresses: %w", config.ErrFeatureDisabled))
	}

	addr, err := svc.repo.GetAddressesByID(ctx, updateModel.UserID, updateModel.ID)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
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

	if err := svc.repo.UpdateUserAddress(ctx, addr); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user address: %w", err))
	}

	addrs, err := svc.repo.GetUserAddresses(ctx, updateModel.UserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load user addresses: %w", err))
	}

	return addrs, nil
}
