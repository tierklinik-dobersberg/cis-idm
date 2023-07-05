package selfservice

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"math/rand"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	"github.com/pquerna/otp/totp"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	idmv1connect.UnimplementedSelfServiceServiceHandler

	cfg  config.Config
	repo *repo.Repo

	common *common.Service
}

func NewService(cfg config.Config, repo *repo.Repo, common *common.Service) (*Service, error) {
	svc := &Service{
		repo:   repo,
		cfg:    cfg,
		common: common,
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

	merr := new(multierror.Error)
	for _, p := range paths {
		switch p {
		case "username":
			if !svc.cfg.FeatureEnabled(config.FeatureAllowUsernameChange) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("username changes are not allowed"))
			}

			user.Username = req.Msg.Username

			if len(user.Username) < 3 {
				merr.Errors = append(merr.Errors, fmt.Errorf("invalid username"))
			}
		case "display_name":
			user.DisplayName = req.Msg.DisplayName
			if len(user.DisplayName) < 3 {
				merr.Errors = append(merr.Errors, fmt.Errorf("invalid display-name"))
			}
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

	mails, err := svc.common.AddEmailAddressToUser(ctx, models.EMail{
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

	mails, err := svc.common.DeleteEmailAddressFromUser(ctx, claims.Subject, req.Msg.Id)
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

	if err := svc.common.MarkEmailAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkEmailAsPrimaryResponse{}), nil
}

func (svc *Service) AddAddress(ctx context.Context, req *connect.Request[idmv1.AddAddressRequest]) (*connect.Response[idmv1.AddAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	addresses, err := svc.common.AddUserAddress(ctx, models.Address{
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

	addresses, err := svc.common.DeleteUserAddress(ctx, claims.Subject, req.Msg.Id)
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

	addrs, err := svc.common.UpdateUserAddress(ctx, models.Address{
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

func (svc *Service) AddPhoneNumber(ctx context.Context, req *connect.Request[idmv1.AddPhoneNumberRequest]) (*connect.Response[idmv1.AddPhoneNumberResponse], error) {
	if !svc.cfg.FeatureEnabled(config.FeaturePhoneNumbers) {
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

	m, err := svc.repo.AddUserPhoneNumber(ctx, m)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.AddPhoneNumberResponse{
		PhoneNumber: conv.PhoneNumberProtoFromPhoneNumber(m),
	}), nil
}

func (svc *Service) DeletePhoneNumber(ctx context.Context, req *connect.Request[idmv1.DeletePhoneNumberRequest]) (*connect.Response[idmv1.DeletePhoneNumberResponse], error) {
	if !svc.cfg.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.repo.DeleteUserPhoneNumber(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.DeletePhoneNumberResponse{}), nil
}

func (svc *Service) MarkPhoneNumberAsPrimary(ctx context.Context, req *connect.Request[idmv1.MarkPhoneNumberAsPrimaryRequest]) (*connect.Response[idmv1.MarkPhoneNumberAsPrimaryResponse], error) {
	if !svc.cfg.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.repo.MarkPhoneNumberAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkPhoneNumberAsPrimaryResponse{}), nil
}

func (svc *Service) Enroll2FA(ctx context.Context, req *connect.Request[idmv1.Enroll2FARequest]) (*connect.Response[idmv1.Enroll2FAResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	switch v := req.Msg.Kind.(type) {
	case *idmv1.Enroll2FARequest_TotpStep1:
		if user.TOTPSecret != "" {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("totp already enrolled"))
		}

		displayName := user.DisplayName
		if displayName == "" {
			displayName = user.Username
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      svc.cfg.SiteName,
			AccountName: displayName,
		})
		if err != nil {
			return nil, err
		}

		mac := hmac.New(sha256.New, []byte(svc.cfg.JWTSecret))

		macString := mac.Sum([]byte(key.Secret()))
		macStringHex := hex.EncodeToString(macString)

		img, err := key.Image(200, 200)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}

		dataUrl := dataurl.EncodeBytes(buf.Bytes())

		return connect.NewResponse(&idmv1.Enroll2FAResponse{
			Kind: &idmv1.Enroll2FAResponse_TotpStep1{
				TotpStep1: &idmv1.EnrollTOTPResponseStep1{
					Secret:     key.Secret(),
					SecretHmac: macStringHex,
					QrCode:     dataUrl,
					Url:        key.String(),
				},
			},
		}), nil

	case *idmv1.Enroll2FARequest_TotpStep2:
		if user.TOTPSecret != "" {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("totp already enrolled"))
		}

		// verify that the secret sent in the request was generated by us
		mac := hmac.New(sha256.New, []byte(svc.cfg.JWTSecret))
		macString := mac.Sum([]byte(v.TotpStep2.Secret))
		macStringHex := hex.EncodeToString(macString)
		if macStringHex != v.TotpStep2.SecretHmac {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid secret"))
		}

		// validate the TOTP passcode
		valid := totp.Validate(v.TotpStep2.VerifyCode, v.TotpStep2.Secret)
		if !valid {
			return nil, fmt.Errorf("invalid passcode")
		}

		if err := svc.repo.SetUserTotpSecret(ctx, claims.Subject, v.TotpStep2.Secret); err != nil {
			return nil, err
		}

		return connect.NewResponse(&idmv1.Enroll2FAResponse{
			Kind: &idmv1.Enroll2FAResponse_TotpStep2{},
		}), nil

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("requested mfa kind is not available"))
	}
}

func (svc *Service) Remove2FA(ctx context.Context, req *connect.Request[idmv1.Remove2FARequest]) (*connect.Response[idmv1.Remove2FAResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	switch v := req.Msg.Kind.(type) {
	case *idmv1.Remove2FARequest_TotpCode:
		if user.TOTPSecret == "" {
			return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("totp 2fa not enrooled"))
		}

		valid := totp.Validate(v.TotpCode, user.TOTPSecret)
		if !valid {
			// check if the user used a recovery code
			recoveryCodeErr := svc.repo.CheckAndDeleteRecoveryCode(ctx, user.ID, v.TotpCode)
			if recoveryCodeErr != nil {
				if errors.Is(recoveryCodeErr, stmts.ErrNoRowsAffected) {
					return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("totp passcode invalid"))
				}

				return nil, recoveryCodeErr
			}
		}

		if err := svc.repo.RemoveUserTotpSecret(ctx, claims.Subject); err != nil {
			return nil, err
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported mfa type"))
	}

	return connect.NewResponse(&idmv1.Remove2FAResponse{}), nil
}

func (svc *Service) GenerateRecoveryCodes(ctx context.Context, req *connect.Request[idmv1.GenerateRecoveryCodesRequest]) (*connect.Response[idmv1.GenerateRecoveryCodesResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)

	codes := make([]string, 20)
	for i := range codes {
		codes[i] = fmt.Sprintf("%d", rand.Intn(999999-100000)+100000)
	}

	if err := svc.repo.ReplaceUserRecoveryCodes(ctx, claims.Subject, codes); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.GenerateRecoveryCodesResponse{
		RecoveryCodes: codes,
	}), nil
}

func (svc *Service) GetRegisteredPasskeys(ctx context.Context, req *connect.Request[idmv1.GetRegisteredPasskeysRequest]) (*connect.Response[idmv1.GetRegisteredPasskeysResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	creds, err := svc.repo.GetPasskeys(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	res := &idmv1.GetRegisteredPasskeysResponse{
		Passkeys: []*idmv1.RegisteredPasskey{},
	}

	for _, cred := range creds {
		res.Passkeys = append(res.Passkeys, &idmv1.RegisteredPasskey{
			Id:           cred.ID,
			ClientName:   cred.ClientName,
			ClientOs:     cred.ClientOS,
			ClientDevice: cred.ClientDevice,
			CredType:     cred.CredType,
		})
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) RemovePasskey(ctx context.Context, req *connect.Request[idmv1.RemovePasskeyRequest]) (*connect.Response[idmv1.RemovePasskeyResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.repo.RemoveWebauthnCred(ctx, claims.Subject, req.Msg.Id); err != nil {
		if errors.Is(err, stmts.ErrNoRowsAffected) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("passkey not found"))
		}

		return nil, err
	}

	return connect.NewResponse(&idmv1.RemovePasskeyResponse{}), nil
}
