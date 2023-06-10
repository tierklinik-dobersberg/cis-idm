package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
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

	validator *protovalidate.Validator
	repo      *repo.Repo
}

func NewService(repo *repo.Repo) (*Service, error) {
	validator, err := protovalidate.New(
		protovalidate.WithMessages(
			new(idmv1.ChangePasswordRequest),
			new(idmv1.UpdateProfileRequest),
			new(idmv1.ValidateEmailRequest),
		),
	)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		repo:      repo,
		validator: validator,
	}

	return svc, nil
}

func (svc *Service) ChangePassword(ctx context.Context, req *connect.Request[idmv1.ChangePasswordRequest]) (*connect.Response[idmv1.ChangePasswordResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no claims associated with request context"))
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get user object: %w", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.GetOldPassword())); err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("incorrect password"))
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Msg.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate password hash: %w", err))
	}

	if err := svc.repo.SetUserPassword(ctx, claims.Subject, string(newHashedPassword)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to save user password: %w", err))
	}

	return connect.NewResponse(&idmv1.ChangePasswordResponse{}), nil
}

func (svc *Service) AddEmailAdress(ctx context.Context, req *connect.Request[idmv1.AddEmailAddressRequest]) (*connect.Response[idmv1.AddEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no token claims associated with request context"))
	}

	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request message: %w", err))
	}

	mails, err := svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get existing user emails: %w", err))
	}

	if _, err := svc.repo.CreateUserEmail(ctx, models.EMail{
		UserID:   claims.Subject,
		Address:  req.Msg.Email,
		Verified: false,
		Primary:  len(mails) == 0, // the first email-address is always marked as primary
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to store new email address: %w", err))
	}

	mails, err = svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get existing user emails: %w", err))
	}

	res := connect.NewResponse(&idmv1.AddEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) DeleteEmailAddress(ctx context.Context, req *connect.Request[idmv1.DeleteEmailAddressRequest]) (*connect.Response[idmv1.DeleteEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no token claims associated with request context"))
	}

	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request message: %w", err))
	}

	if err := svc.repo.DeleteEMailFromUser(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete email from user: %w", err))
	}

	mails, err := svc.repo.GetUserEmails(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get existing user emails: %w", err))
	}

	res := connect.NewResponse(&idmv1.DeleteEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}
