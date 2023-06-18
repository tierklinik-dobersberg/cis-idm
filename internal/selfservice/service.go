package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	idmv1connect.UnimplementedSelfServiceServiceHandler

	repo *repo.Repo
}

func NewService(repo *repo.Repo) (*Service, error) {
	svc := &Service{
		repo: repo,
	}

	return svc, nil
}

func (svc *Service) UpdateProfile(ctx context.Context, req *connect.Request[idmv1.UpdateProfileRequest]) (*connect.Response[idmv1.UpdateProfileResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get user object: %w", err)
	}

	paths := req.Msg.GetFieldMask().GetPaths()
	if len(paths) == 0 {
		paths = []string{"username", "display_name", "first_name", "last_name", "avatar", "birthday"}
	}

	for _, p := range paths {
		switch p {
		case "username":
			user.Username = req.Msg.Username
		case "display_name":
			user.DisplayName = req.Msg.DisplayName
		case "first_name":
			user.FirstName = req.Msg.FirstName
		case "last_name":
			user.LastName = req.Msg.LastName
		case "avatar":
			user.Avatar = req.Msg.Avatar
		case "birthday":
			user.Birthday = req.Msg.Birthday
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid field mask for update operation: invalid field name %q", p))
		}
	}

	merr := new(multierror.Error)
	if len(user.Username) < 3 {
		merr.Errors = append(merr.Errors, fmt.Errorf("invalid username"))
	}
	if len(user.DisplayName) < 3 {
		merr.Errors = append(merr.Errors, fmt.Errorf("invalid display-name"))
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := svc.repo.UpdateUser(ctx, user); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user: %w", err))
	}

	return connect.NewResponse(&idmv1.UpdateProfileResponse{
		User: conv.UserProtoFromUser(user),
	}), nil

}

func (svc *Service) ChangePassword(ctx context.Context, req *connect.Request[idmv1.ChangePasswordRequest]) (*connect.Response[idmv1.ChangePasswordResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get user object: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.GetOldPassword())); err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("incorrect password"))
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Msg.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	if err := svc.repo.SetUserPassword(ctx, claims.Subject, string(newHashedPassword)); err != nil {
		return nil, fmt.Errorf("failed to save user password: %w", err)
	}

	return connect.NewResponse(&idmv1.ChangePasswordResponse{}), nil
}

func (svc *Service) AddEmailAddress(ctx context.Context, req *connect.Request[idmv1.AddEmailAddressRequest]) (*connect.Response[idmv1.AddEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	mails, err := svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	if _, err := svc.repo.CreateUserEmail(ctx, models.EMail{
		UserID:   claims.Subject,
		Address:  req.Msg.Email,
		Verified: false,
		Primary:  len(mails) == 0, // the first email-address is always marked as primary
	}); err != nil {
		return nil, fmt.Errorf("failed to store new email address: %w", err)
	}

	mails, err = svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user emails: %w", err)
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

	middleware.L(ctx).WithField("email_id", req.Msg.Id).Infof("deleting email address from user")

	if err := svc.repo.DeleteEMailFromUser(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, fmt.Errorf("failed to delete email from user: %w", err)
	}

	mails, err := svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	res := connect.NewResponse(&idmv1.DeleteEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) AddAddress(ctx context.Context, req *connect.Request[idmv1.AddAddressRequest]) (*connect.Response[idmv1.AddAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if _, err := svc.repo.AddUserAddress(ctx, models.Address{
		UserID:   claims.Subject,
		CityCode: req.Msg.CityCode,
		CityName: req.Msg.CityName,
		Street:   req.Msg.Street,
		Extra:    req.Msg.Extra,
	}); err != nil {
		return nil, fmt.Errorf("failed to save new user address: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
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

	if err := svc.repo.DeleteUserAddress(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, fmt.Errorf("failed to delete user address: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
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

	addr, err := svc.repo.GetAddressesByID(ctx, claims.Subject, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load address by id: %w", err))
	}

	paths := req.Msg.GetFieldMask().GetPaths()
	if len(paths) == 0 {
		paths = []string{
			"city_code", "city", "street", "extra",
		}
	}

	for _, p := range paths {
		switch p {
		case "city_code":
			addr.CityCode = req.Msg.CityCode
		case "city_name":
			addr.CityName = req.Msg.CityName
		case "street":
			addr.Street = req.Msg.Street
		case "extra":
			addr.Extra = req.Msg.Extra
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid field_mask for update operation: invalid path %q", p))
		}
	}

	if err := svc.repo.UpdateUserAddress(ctx, addr); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user address: %w", err))
	}

	addrs, err := svc.repo.GetUserAddresses(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to load user addresses: %w", err))
	}

	return connect.NewResponse(&idmv1.UpdateAddressResponse{
		Addresses: conv.AddressProtosFromAddresses(addrs...),
	}), nil
}
