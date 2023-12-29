package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/mennanov/fmutils"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/mailer"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type AuthService struct {
	idmv1connect.UnimplementedAuthServiceHandler

	*app.Providers
}

// NewService returns a new authentication service that verifies users using repo.
func NewService(providers *app.Providers) *AuthService {
	return &AuthService{
		Providers: providers,
	}
}

func (svc *AuthService) Login(ctx context.Context, req *connect.Request[idmv1.LoginRequest]) (*connect.Response[idmv1.LoginResponse], error) {
	r := req.Msg

	var (
		user repo.User
	)

	// Log out any user that might still be logged-in.
	claims := middleware.ClaimsFromContext(ctx)
	if claims != nil {
		// There's already a user logged in so this seems like a forced user-switch.
		// In this case, we should invalidate the current token (and refresh token)
		// as well as any web-push subscriptions
		if err := svc.invalidateTokens(ctx, claims); err != nil {
			return nil, err
		}
	}

	kind := "password"

	switch r.AuthType {
	case idmv1.AuthType_AUTH_TYPE_PASSWORD:
		if r.AuthType != idmv1.AuthType_AUTH_TYPE_PASSWORD {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("auth type not supported"))
		}

		passwordAuth := r.GetPassword()
		if passwordAuth == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid payload for password auth type"))
		}

		logrus.Infof("authentication request for user %s", passwordAuth.GetUsername())

		var err error
		user, err = svc.Datastore.GetUserByName(ctx, passwordAuth.GetUsername())
		if err != nil {
			if svc.Config.FeatureEnabled(config.FeatureLoginByMail) {
				if errors.Is(err, sql.ErrNoRows) {
					response, err := svc.Datastore.GetUserByEMail(ctx, passwordAuth.GetUsername())

					user = response.User

					if err == nil && !response.Verified {
						return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("e-mail address has not been verified"))
					}
				}
			}

			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not found: %w", err))
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordAuth.GetPassword())); err != nil {
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("incorrect password"))
			}
		}

		// check if the user still needs to pass the 2fa
		if user.TotpSecret.String != "" {
			state, _, err := svc.CreateSignedJWT(user, nil, "", time.Minute*5, jwt.Scope2FAPending)
			if err != nil {
				return nil, err
			}

			return connect.NewResponse(&idmv1.LoginResponse{
				Response: &idmv1.LoginResponse_MfaRequired{
					MfaRequired: &idmv1.MFARequiredResponse{
						Kind:  idmv1.RequiredMFAKind_REQUIRED_MFA_KIND_TOTP,
						State: state,
					},
				},
			}), nil
		}
		// otherwise continue outside of the switch block and issue access and refresh tokens

	case idmv1.AuthType_AUTH_TYPE_TOTP:
		if req.Msg.GetTotp() == nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid message"))
		}

		claims, err := jwt.ParseAndVerify([]byte(svc.Config.JWTSecret), req.Msg.GetTotp().State)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		if !slices.Contains(claims.Scopes, jwt.Scope2FAPending) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid message"))
		}

		user, err = svc.Datastore.GetUserByID(ctx, claims.Subject)
		if err != nil {
			return nil, err
		}

		if user.TotpSecret.String == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("totp not enrolled"))
		}

		valid := totp.Validate(req.Msg.GetTotp().Code, user.TotpSecret.String)
		if !valid {
			// if the code is not valid the user might used a recovery code.
			// TODO(ppacher): do we have security implications if we automatically try
			// recovery codes here?
			rows, recoveryCodeErr := svc.Datastore.CheckAndDeleteRecoveryCode(ctx, repo.CheckAndDeleteRecoveryCodeParams{
				UserID: claims.Subject,
				Code:   req.Msg.GetTotp().Code,
			})

			if recoveryCodeErr != nil {

				// any other internal error
				return nil, err
			}

			if rows == 0 {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid totp passcode"))
			}
		}

		kind = "mfa"

		// continue outside of the switch block and issue access and refresh tokens
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported authentication method"))
	}

	roles, err := svc.Datastore.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// make sure we provide a display name in the response.
	// - either join first and lastname
	// - or fall back to the user name.
	common.EnsureDisplayName(&user)

	response := &idmv1.AccessTokenResponse{
		User: &idmv1.User{
			Id:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
		},
	}

	redirectTo, err := svc.HandleRequestedRedirect(ctx, req.Msg.RequestedRedirect)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&idmv1.LoginResponse{
		Response: &idmv1.LoginResponse_AccessToken{
			AccessToken: response,
		},
		RedirectTo: redirectTo,
	})

	var refreshTokenID string

	if !r.GetNoRefreshToken() {
		_, refreshTokenID, err = svc.AddRefreshToken(user, roles, kind, resp.Header())
		if err != nil {
			return nil, err
		}
	}

	if token, _, err := svc.AddAccessToken(user, roles, req.Msg.Ttl.AsDuration(), refreshTokenID, kind, resp.Header()); err != nil {
		return nil, err
	} else {
		response.Token = token
	}

	return resp, nil
}

func (svc *AuthService) Logout(ctx context.Context, req *connect.Request[idmv1.LogoutRequest]) (*connect.Response[idmv1.LogoutResponse], error) {
	// get the JWT token claims from the request context
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	if err := svc.invalidateTokens(ctx, claims); err != nil {
		return nil, err
	}

	resp := connect.NewResponse(new(idmv1.LogoutResponse))

	// clear the refresh token cookie
	clearRefreshCookie := http.Cookie{
		Name:     svc.Config.RefreshTokenCookieName,
		Value:    "",
		Domain:   svc.Config.Domain,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Path:     "/tkd.idm.v1.AuthService/RefreshToken",
		HttpOnly: true,
	}

	clearAccessCookie := http.Cookie{
		Name:     svc.Config.AccessTokenCookieName,
		Value:    "",
		Domain:   svc.Config.Domain,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		HttpOnly: true,
	}

	resp.Header().Add("Set-Cookie", clearRefreshCookie.String())
	resp.Header().Add("Set-Cookie", clearAccessCookie.String())
	resp.Header().Add("Clear-Site-Data", `"cookies"`)

	return resp, nil
}

func (svc *AuthService) RefreshToken(ctx context.Context, req *connect.Request[idmv1.RefreshTokenRequest]) (*connect.Response[idmv1.RefreshTokenResponse], error) {
	refreshCookie := middleware.FindCookie(svc.Config.RefreshTokenCookieName, req.Header())
	if refreshCookie == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no refresh cookie provided"))
	}

	claims, err := jwt.ParseAndVerify([]byte(svc.Config.JWTSecret), refreshCookie.Value)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid refresh token: %w", err))
	}

	if !slices.Contains(claims.Scopes, jwt.ScopeRefresh) {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("the provided token is invalid"))
	}

	user, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid refresh token"))
	}

	roles, err := svc.Datastore.GetRolesForUser(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to get role assignments: %w", err)
	}

	redirectTo, err := svc.HandleRequestedRedirect(ctx, req.Msg.RequestedRedirect)
	if err != nil {
		return nil, err
	}

	tokenResponse := &idmv1.AccessTokenResponse{
		User: conv.UserProtoFromUser(ctx, user),
	}

	resp := connect.NewResponse(&idmv1.RefreshTokenResponse{
		AccessToken: tokenResponse,
		RedirectTo:  redirectTo,
	})

	kind := ""
	if claims.AppMetadata != nil {
		kind = claims.AppMetadata.LoginKind
	}
	token, _, err := svc.AddAccessToken(user, roles, req.Msg.Ttl.AsDuration(), claims.ID, kind, resp.Header())
	if err != nil {
		return nil, err
	}

	tokenResponse.Token = token

	return resp, nil
}

func (svc *AuthService) Introspect(ctx context.Context, req *connect.Request[idmv1.IntrospectRequest]) (*connect.Response[idmv1.IntrospectResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("not jwt token claims found in request context")
	}

	user, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user")
	}

	profile, err := svc.GetUserProfileProto(ctx, user)
	if err != nil {
		return nil, err
	}

	if paths := req.Msg.GetReadMask().GetPaths(); len(paths) > 0 {
		if req.Msg.ExcludeFields {
			fmutils.Prune(profile, paths)
		} else {
			fmutils.Filter(profile, paths)
		}
	}

	return connect.NewResponse(&idmv1.IntrospectResponse{
		Profile: profile,
	}), nil
}

func (svc *AuthService) GenerateRegistrationToken(ctx context.Context, req *connect.Request[idmv1.GenerateRegistrationTokenRequest]) (*connect.Response[idmv1.GenerateRegistrationTokenResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("not jwt token claims found in request context")
	}

	creator, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid jwt token or account deleted")
		}

		return nil, err
	}

	token, err := svc.Providers.GenerateRegistrationToken(ctx, creator, req.Msg.MaxCount, req.Msg.GetTtl().AsDuration(), req.Msg.InitialRoles)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.GenerateRegistrationTokenResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) ValidateRegistrationToken(ctx context.Context, req *connect.Request[idmv1.ValidateRegistrationTokenRequest]) (*connect.Response[idmv1.ValidateRegistrationTokenResponse], error) {
	_, err := svc.Datastore.ValidateRegistrationToken(ctx, req.Msg.Token)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.ValidateRegistrationTokenResponse{}), nil
}

func (svc *AuthService) RegisterUser(ctx context.Context, req *connect.Request[idmv1.RegisterUserRequest]) (*connect.Response[idmv1.RegisterUserResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeatureSelfRegistration) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("registration feature is disabled"))
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := svc.Datastore.Tx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.L(ctx).Errorf("failed to rollback transaction: %s", err)
		}
	}()

	userModel, err := svc.CreateUser(ctx, tx, repo.CreateUserParams{
		Username: req.Msg.Username,
		Password: string(passwordHash),
	}, req.Msg.RegistrationToken)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	if mailModel, err := svc.Datastore.
		WithTx(tx).
		CreateEMail(ctx, repo.CreateEMailParams{
		ID: id.String(),
		UserID:    userModel.ID,
		Address:   req.Msg.Email,
		IsPrimary: true,
	}); err != nil {
		// just log out the error but continue to sign in the user
		log.L(ctx).WithError(err).Errorf("failed to save email address for user")
	} else {
		if err := svc.SendMailVerification(ctx, *userModel, mailModel); err != nil {
			log.L(ctx).WithError(err).Errorf("failed to send verification mail")
		}
	}

	if err := tx.Commit(); err != nil {
		log.L(ctx).Errorf("failed to commit transaction: %s", err)

		return nil, err
	}

	roles, err := svc.Datastore.GetRolesForUser(ctx, userModel.ID)
	if err != nil {
		log.L(ctx).WithError(err).Error("failed to get user role assignments")
	}

	tokenResponse := &idmv1.AccessTokenResponse{
		User: conv.UserProtoFromUser(ctx, *userModel),
	}

	resp := connect.NewResponse(&idmv1.RegisterUserResponse{
		AccessToken: tokenResponse,
	})

	_, refreshTokenID, err := svc.AddRefreshToken(*userModel, roles, "password", resp.Header())
	if err != nil {
		return nil, err
	}

	token, _, err := svc.AddAccessToken(*userModel, roles, 0, refreshTokenID, "password", resp.Header())
	if err != nil {
		return nil, err
	}

	tokenResponse.Token = token

	return resp, nil
}

func (svc *AuthService) CreateUser(ctx context.Context, tx *sql.Tx, params repo.CreateUserParams, token string) (*repo.User, error) {
	var (
		initialRoles []string
		err          error
	)

	db := svc.Datastore.WithTx(tx)

	if params.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		params.ID = id.String()
	}

	// ensure we have a valid registration token if IDM_REGISTRATION_REQUIRES_TOKEN is set to true.
	// Note that we also accept a registration token even if it's not required so users can be
	// bootstrapped with a set of initial roles.
	if svc.Config.RegistrationRequiresToken || token != "" {
		tokenModel, err := db.GetRegistrationToken(ctx, repo.GetRegistrationTokenParams{
			Token: token,
			Expires: sql.NullTime{
				Time: time.Now(),
			},
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid registration token"))
			}
			return nil, err
		}

		if len(tokenModel.InitialRoles) > 0 {
			if err := json.Unmarshal([]byte(tokenModel.InitialRoles), &initialRoles); err != nil {
				return nil, err
			}
		}

		if _, err := db.MarkRegistrationTokenUsed(ctx, repo.MarkRegistrationTokenUsedParams{
			Token: token,
			Expires: sql.NullTime{
				Time: time.Now(),
			},
		}); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid registration token"))
			}
			return nil, err
		}
	}

	userModel, err := db.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	merr := new(multierror.Error)
	for _, role := range initialRoles {
		if err := db.AssignRoleToUser(ctx, repo.AssignRoleToUserParams{
			UserID: userModel.ID,
			RoleID: role,
		}); err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to assign role %s: %w", role, err))
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return &userModel, err
	}

	return &userModel, nil
}

func (svc *AuthService) RequestPasswordReset(ctx context.Context, req *connect.Request[idmv1.RequestPasswordResetRequest]) (*connect.Response[idmv1.RequestPasswordResetResponse], error) {
	switch v := req.Msg.Kind.(type) {
	case *idmv1.RequestPasswordResetRequest_Email:
		if v.Email == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("missing username or email address"))
		}

		var user repo.User
		response, err := svc.Datastore.GetUserByEMail(ctx, v.Email)
		if err == nil {
			user = response.User
		} else {
			user, err = svc.Datastore.GetUserByName(ctx, v.Email)
			if err != nil {
				return connect.NewResponse(&idmv1.RequestPasswordResetResponse{}), nil
			}

			primaryMail, err := svc.Datastore.GetPrimaryEmailForUserByID(ctx, user.ID)
			if err != nil {
				return connect.NewResponse(&idmv1.RequestPasswordResetResponse{}), nil
			}

			v.Email = primaryMail.Address
		}

		code, cacheKey, err := svc.Common.GeneratePasswordResetToken(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		// Send a text message to the user
		msg := mailer.Message{
			From: svc.Config.MailConfig.From,
			To:   []string{v.Email},
		}

		// make sure we have a valid display name for the user
		common.EnsureDisplayName(&user)

		if err := mailer.SendTemplate(ctx, svc.Config, svc.TemplateEngine, svc.Mailer, msg, tmpl.RequestPasswordReset, &tmpl.RequestPasswordResetCtx{
			User:      user,
			ResetLink: fmt.Sprintf(svc.Config.PasswordResetURL, code),
		}); err != nil {
			defer func() {
				_ = svc.Cache.DeleteKey(ctx, cacheKey)
			}()

			return nil, err
		}

	case *idmv1.RequestPasswordResetRequest_PasswordReset:
		cacheKey := fmt.Sprintf("password-reset:%s", v.PasswordReset.Token)
		var userID string
		if err := svc.Cache.GetAndDeleteKey(ctx, cacheKey, &userID); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid or expired token"))
		}

		user, err := svc.Datastore.GetUserByID(ctx, userID)
		if err != nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(v.PasswordReset.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		rows, err := svc.Datastore.SetUserPassword(ctx, repo.SetUserPasswordParams{
			Password: string(hashed),
			ID:       user.ID,
		})
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request"))
	}

	return connect.NewResponse(&idmv1.RequestPasswordResetResponse{}), nil
}

func (svc *AuthService) invalidateTokens(ctx context.Context, claims *jwt.Claims) error {
	// FIXME(ppacher): do not abort if access-token invalidation fails

	// mark the token as rejected
	if err := common.Timing(ctx, "reject-access-token", func() error {
		// delete the web-push subscription for the access token
		_, _ = svc.Datastore.DeleteWebPushSubscriptionForToken(ctx, claims.ID)

		return svc.Datastore.CreateRejectedToken(ctx, repo.CreateRejectedTokenParams{
			TokenID:   claims.ID,
			UserID:    claims.Subject,
			IssuedAt:  time.Unix(claims.IssuedAt, 0),
			ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		})
	}); err != nil {
		return fmt.Errorf("failed to mark token as rejected: %w", err)
	}

	// also mark the parent (refresh) token as rejected
	if claims.AppMetadata != nil && claims.AppMetadata.ParentTokenID != "" {
		if err := common.Timing(ctx, "reject-refresh-token", func() error {
			// delete the web-push subscription for the refresh token
			_, _ = svc.Datastore.DeleteWebPushSubscriptionForToken(ctx, claims.AppMetadata.ParentTokenID)

			return svc.Datastore.CreateRejectedToken(ctx, repo.CreateRejectedTokenParams{
				TokenID:   claims.AppMetadata.ParentTokenID,
				UserID:    claims.Subject,
				IssuedAt:  time.Unix(claims.IssuedAt, 0),
				ExpiresAt: time.Unix(claims.ExpiresAt, 0),
			})
		}); err != nil {
			return fmt.Errorf("failed to mark token as rejected: %w", err)
		}
	}

	return nil
}

var _ idmv1connect.AuthServiceHandler = new(AuthService)
