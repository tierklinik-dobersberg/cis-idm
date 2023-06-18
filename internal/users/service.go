package users

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	"github.com/mennanov/fmutils"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	idmv1connect.UnimplementedUserServiceHandler

	repo *repo.Repo
}

func NewService(repo *repo.Repo) (*Service, error) {
	svc := &Service{
		repo: repo,
	}

	return svc, nil
}

func (svc *Service) ListUsers(ctx context.Context, req *connect.Request[idmv1.ListUsersRequest]) (*connect.Response[idmv1.ListUsersResponse], error) {
	users, err := svc.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &idmv1.ListUsersResponse{}
	for _, usr := range users {
		addresses, err := svc.repo.GetUserAddresses(ctx, usr.ID)
		if err != nil {
			middleware.L(ctx).Errorf("failed to get user addresses for user %s: %s", usr.ID, err)
		}

		mails, err := svc.repo.GetUserEmails(ctx, usr.ID)
		if err != nil {
			middleware.L(ctx).Errorf("failed to get user emails for user %s: %s", usr.ID, err)
		}

		phones, err := svc.repo.GetUserPhoneNumbers(ctx, usr.ID)
		if err != nil {
			middleware.L(ctx).Errorf("failed to get user phone numbers for user %s: %s", usr.ID, err)
		}

		profileProto := conv.ProfileProtoFromUser(
			usr,
			conv.WithAddresses(addresses...),
			conv.WithEmailAddresses(mails...),
			conv.WithPhoneNumbers(phones...),
		)

		res.Users = append(res.Users, profileProto)
	}

	// make sure we only include fields that are requested
	if fieldMaskPaths := req.Msg.GetFieldMask().GetPaths(); len(fieldMaskPaths) > 0 {
		fmutils.Filter(res, fieldMaskPaths)
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) GetUser(ctx context.Context, req *connect.Request[idmv1.GetUserRequest]) (*connect.Response[idmv1.GetUserResponse], error) {
	user, err := svc.repo.GetUserByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.GetUserResponse{
		User: conv.UserProtoFromUser(user),
	}), nil
}

func (svc *Service) DeleteUser(ctx context.Context, req *connect.Request[idmv1.DeleteUserRequest]) (*connect.Response[idmv1.DeleteUserResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no jwt claims associated with request context"))
	}

	// make sure users just cannot directly delete their own profile.
	// TODO(ppacher): add support to "request profile deletion"
	if claims.Subject == req.Msg.Id {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("deleting your own account is not allowed"))
	}

	// actually delete the user from the repository
	if err := svc.repo.DeleteUser(ctx, req.Msg.Id); err != nil {
		return nil, err
	}

	// TODO(ppacher): invalidate all current access and refresh tokens

	return connect.NewResponse(&idmv1.DeleteUserResponse{}), nil
}

func (svc *Service) CreateUser(ctx context.Context, req *connect.Request[idmv1.CreateUserRequest]) (*connect.Response[idmv1.CreateUserResponse], error) {
	if req.Msg.Profile == nil || req.Msg.Profile.User == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request: missing user object"))
	}
	usr := req.Msg.Profile.User

	userModel := models.User{
		Username:    usr.Username,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		DisplayName: usr.DisplayName,
		Avatar:      usr.Avatar,
		Birthday:    usr.Birthday,
	}

	if req.Msg.Password != "" {
		if req.Msg.PasswordIsBcrypt {
			userModel.Password = req.Msg.Password
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}

			userModel.Password = string(hash)
		}
	}

	// actually create the user.
	userModel, err := svc.repo.CreateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	merr := new(multierror.Error)

	// Create user addresses
	var userAddresses []models.Address
	if addresses := req.Msg.GetProfile().Addresses; len(addresses) > 0 {
		for _, addr := range addresses {
			addrModel := models.Address{
				UserID:   userModel.ID,
				CityCode: addr.CityCode,
				CityName: addr.CityName,
				Street:   addr.Street,
				Extra:    addr.Extra,
			}

			if addr, err := svc.repo.AddUserAddress(ctx, addrModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create user address: %w", err))
			} else {
				userAddresses = append(userAddresses, addr)
			}
		}
	}

	// Create phone number records
	var userPhoneNumbers []models.PhoneNumber
	if phoneNumbers := req.Msg.GetProfile().PhoneNumbers; len(phoneNumbers) > 0 {
		for idx, nbr := range phoneNumbers {
			nbrModel := models.PhoneNumber{
				UserID:      userModel.ID,
				PhoneNumber: nbr,
				Primary:     idx == 0,
			}

			if phone, err := svc.repo.AddUserPhoneNumber(ctx, nbrModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create phone number: %w", err))
			} else {
				userPhoneNumbers = append(userPhoneNumbers, phone)
			}
		}
	}

	// create email address records
	var userEmails []models.EMail
	if emails := req.Msg.GetProfile().EmailAddresses; len(emails) > 0 {
		for idx, mail := range emails {
			mailModel := models.EMail{
				UserID:   userModel.ID,
				Address:  mail.Address,
				Verified: true,
				Primary:  idx == 0,
			}

			if email, err := svc.repo.CreateUserEmail(ctx, mailModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create email record: %w", err))
			} else {
				userEmails = append(userEmails, email)
			}
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&idmv1.CreateUserResponse{
		Profile: conv.ProfileProtoFromUser(
			userModel,
			conv.WithAddresses(userAddresses...),
			conv.WithEmailAddresses(userEmails...),
			conv.WithPhoneNumbers(userPhoneNumbers...),
		),
	}), nil
}
