package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type AuthService struct {
	idmv1connect.UnimplementedAuthServiceHandler

	validator *protovalidate.Validator

	repo *repo.Repo
	cfg  config.Config
}

// NewService returns a new authentication service that verifies users using repo.
func NewService(repo *repo.Repo, cfg config.Config) (*AuthService, error) {
	validator, err := protovalidate.New(
		protovalidate.WithMessages(
			&idmv1.LoginRequest{},
			&idmv1.LogoutRequest{},
			&idmv1.IntrospectRequest{},
			&idmv1.RefreshTokenRequest{},
		),
	)

	if err != nil {
		return nil, err
	}

	return &AuthService{
		repo:      repo,
		cfg:       cfg,
		validator: validator,
	}, nil
}

func (svc *AuthService) Login(ctx context.Context, req *connect.Request[idmv1.LoginRequest]) (*connect.Response[idmv1.LoginResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	logrus.Infof("received authentication request")
	r := req.Msg

	var (
		user models.User
	)

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
		user, err = svc.repo.GetUserByName(ctx, passwordAuth.GetUsername())
		if err != nil {
			if svc.cfg.FeatureEnabled(config.FeatureLoginByMail) {
				if errors.Is(err, stmts.ErrNoResults) {
					var verified bool
					user, verified, err = svc.repo.GetUserByEMail(ctx, passwordAuth.GetUsername())

					if err == nil && !verified {
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
		if user.TOTPSecret != "" {
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

		claims, err := jwt.ParseAndVerify([]byte(svc.cfg.JWTSecret), req.Msg.GetTotp().State)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		if !slices.Contains(claims.Scopes, jwt.Scope2FAPending) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid message"))
		}

		user, err = svc.repo.GetUserByID(ctx, claims.Subject)
		if err != nil {
			return nil, err
		}

		if user.TOTPSecret == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("totp not enrolled"))
		}

		valid := totp.Validate(req.Msg.GetTotp().Code, user.TOTPSecret)
		if !valid {
			// if the code is not valid the user might used a recovery code.
			// TODO(ppacher): do we have security implications if we automatically try
			// recovery codes here?
			recoveryCodeErr := svc.repo.CheckAndDeleteRecoveryCode(ctx, claims.Subject, req.Msg.GetTotp().Code)
			if recoveryCodeErr != nil {
				if errors.Is(recoveryCodeErr, stmts.ErrNoRowsAffected) {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid totp passcode"))
				}

				// any other internal error
				return nil, err
			}
		}

		// continue outside of the switch block and issue access and refresh tokens
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported authentication method"))
	}

	roles, err := svc.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// make sure we provide a display name in the response.
	// - either join first and lastname
	// - or fall back to the user name.
	if user.DisplayName == "" {
		switch {
		case user.FirstName != "" || user.LastName != "":
			user.DisplayName = strings.Join([]string{user.FirstName, user.LastName}, " ")
		default:
			user.DisplayName = user.Username
		}
	}

	var (
		refreshToken   string
		refreshTokenID string
	)

	if !r.GetNoRefreshToken() {
		refreshToken, refreshTokenID, err = svc.CreateSignedJWT(user, nil, "", svc.cfg.RefreshTokenTTL.AsDuration(), jwt.ScopeRefresh)
		if err != nil {
			return nil, err
		}
	}

	accessTokenTTL := svc.cfg.AccessTokenTTL.AsDuration()
	if req.Msg.Ttl != nil {
		requestedTTL := req.Msg.Ttl.AsDuration()
		if requestedTTL > 0 && requestedTTL < accessTokenTTL {
			accessTokenTTL = requestedTTL
		}
	}

	accessToken, _, err := svc.CreateSignedJWT(user, roles, refreshTokenID, accessTokenTTL, jwt.ScopeAccess)
	if err != nil {
		return nil, err
	}

	redirectTo, err := svc.HandleRequestedRedirect(ctx, req.Msg.RequestedRedirect)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&idmv1.LoginResponse{
		Response: &idmv1.LoginResponse_AccessToken{
			AccessToken: &idmv1.AccessTokenResponse{
				Token: accessToken,
				User: &idmv1.User{
					Id:          user.ID,
					Username:    user.Username,
					DisplayName: user.DisplayName,
				},
			},
		},
		RedirectTo: redirectTo,
	})

	svc.AddAccessTokenCookie(resp.Header(), accessToken, accessTokenTTL)

	if !r.GetNoRefreshToken() {
		svc.AddRefreshTokenCookie(resp.Header(), refreshToken, svc.cfg.RefreshTokenTTL.AsDuration())
	}

	return resp, nil
}

func (svc *AuthService) Logout(ctx context.Context, req *connect.Request[idmv1.LogoutRequest]) (*connect.Response[idmv1.LogoutResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// get the JWT token claims from the request context
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	// mark the token as rejected
	if err := svc.repo.MarkTokenRejected(ctx, models.RejectedToken{
		TokenID:   claims.ID,
		UserID:    claims.Subject,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}); err != nil {
		return nil, fmt.Errorf("failed to mark token as rejected: %w", err)
	}

	// also mark the parent (refresh) token as rejected
	if claims.AppMetadata != nil && claims.AppMetadata.ParentTokenID != "" {
		if err := svc.repo.MarkTokenRejected(ctx, models.RejectedToken{
			TokenID:   claims.AppMetadata.ParentTokenID,
			UserID:    claims.Subject,
			IssuedAt:  claims.IssuedAt,
			ExpiresAt: claims.ExpiresAt,
		}); err != nil {
			return nil, fmt.Errorf("failed to mark token as rejected: %w", err)
		}
	}

	resp := connect.NewResponse(new(idmv1.LogoutResponse))

	// clear the refresh token cookie
	clearRefreshCookie := http.Cookie{
		Name:     svc.cfg.RefreshTokenCookieName,
		Value:    "",
		Domain:   svc.cfg.Domain,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Path:     "/tkd.idm.v1.AuthService/RefreshToken",
		HttpOnly: true,
	}

	clearAccessCookie := http.Cookie{
		Name:     svc.cfg.AccessTokenCookieName,
		Value:    "",
		Domain:   svc.cfg.Domain,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		HttpOnly: true,
	}

	resp.Header().Add("Set-Cookie", clearRefreshCookie.String())
	resp.Header().Add("Set-Cookie", clearAccessCookie.String())
	resp.Header().Add("Clear-Site-Data", `"cache", "cookies"`) // we keep localStorage for the loggedInUsers key

	return resp, nil
}

func (svc *AuthService) RefreshToken(ctx context.Context, req *connect.Request[idmv1.RefreshTokenRequest]) (*connect.Response[idmv1.RefreshTokenResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	refreshCookie := middleware.FindCookie(svc.cfg.RefreshTokenCookieName, req.Header())
	if refreshCookie == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no refresh cookie provided"))
	}

	claims, err := jwt.ParseAndVerify([]byte(svc.cfg.JWTSecret), refreshCookie.Value)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid refresh token: %w", err))
	}

	if !slices.Contains(claims.Scopes, jwt.ScopeRefresh) {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("the provided token is invalid"))
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid refresh token"))
	}

	roles, err := svc.repo.GetUserRoles(ctx, claims.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get group memberships: %w", err)
	}

	tokenTTL := svc.cfg.AccessTokenTTL.AsDuration()
	token, _, err := svc.CreateSignedJWT(user, roles, claims.ID, tokenTTL, jwt.ScopeAccess)
	if err != nil {
		return nil, err
	}

	redirectTo, err := svc.HandleRequestedRedirect(ctx, req.Msg.RequestedRedirect)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&idmv1.RefreshTokenResponse{
		AccessToken: &idmv1.AccessTokenResponse{
			Token: token,
			User:  conv.UserProtoFromUser(user),
		},
		RedirectTo: redirectTo,
	})

	svc.AddAccessTokenCookie(resp.Header(), token, tokenTTL)

	return resp, nil
}

func (svc *AuthService) Introspect(ctx context.Context, req *connect.Request[idmv1.IntrospectRequest]) (*connect.Response[idmv1.IntrospectResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("not jwt token claims found in request context")
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid user")
	}

	emails, err := svc.repo.GetUserEmails(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load emails: %w", err)
	}

	addresses, err := svc.repo.GetUserAddresses(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user addresses: %w", err)
	}

	phoneNumbers, err := svc.repo.GetUserPhoneNumbers(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user phone numbers: %w", err)
	}

	var primaryMail *models.EMail
	for _, mail := range emails {
		if mail.Primary {
			primaryMail = new(models.EMail)
			*primaryMail = mail

			break
		}
	}

	var primaryPhone *models.PhoneNumber
	for _, phone := range phoneNumbers {
		if phone.Primary {
			primaryPhone = new(models.PhoneNumber)
			*primaryPhone = phone

			break
		}
	}

	hasBackupCodes, err := svc.repo.UserHasRecoveryCodes(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing recovery codes: %w", err)
	}

	return connect.NewResponse(&idmv1.IntrospectResponse{
		Profile: conv.ProfileProtoFromUser(
			user,
			conv.WithUserHasRecoveryCodes(hasBackupCodes),
			conv.WithEmailAddresses(emails...),
			conv.WithAddresses(addresses...),
			conv.WithPhoneNumbers(phoneNumbers...),
			conv.WithPrimaryMail(primaryMail),
			conv.WithPrimaryPhone(primaryPhone),
		),
	}), nil
}

func (svc *AuthService) GenerateRegistrationToken(ctx context.Context, req *connect.Request[idmv1.GenerateRegistrationTokenRequest]) (*connect.Response[idmv1.GenerateRegistrationTokenResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("not jwt token claims found in request context")
	}

	token, err := bootstrap.GenerateSecret(8)
	if err != nil {
		return nil, err
	}

	tokenModel := models.RegistrationToken{
		Token:     token,
		CreatedBy: claims.Subject,
		CreatedAt: time.Now().Unix(),
	}

	if req.Msg.MaxCount > 0 {
		i := new(int64)
		*i = int64(req.Msg.MaxCount)
		tokenModel.AllowedUsage = i
	}

	if req.Msg.Ttl != nil {
		expires := time.Now().Add(req.Msg.Ttl.AsDuration()).Unix()
		tokenModel.Expires = &expires
	}

	if len(req.Msg.InitialRoles) > 0 {
		var initialRoles []string

		for _, role := range req.Msg.InitialRoles {
			roleModel, err := svc.repo.GetRoleByID(ctx, role)
			if err != nil {
				if errors.Is(err, stmts.ErrNoResults) {
					roleModel, err = svc.repo.GetRoleByName(ctx, role)
				}
			}

			if err != nil {
				if errors.Is(err, stmts.ErrNoResults) {
					return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("roles %q does not exist", role))
				}

				return nil, err
			}
			initialRoles = append(initialRoles, roleModel.ID)
		}

		roleBlob, err := json.Marshal(initialRoles)
		if err != nil {
			return nil, err
		}

		tokenModel.InitialRoles = string(roleBlob)
	}

	if err := svc.repo.CreateRegistrationToken(ctx, tokenModel); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.GenerateRegistrationTokenResponse{
		Token: token,
	}), nil
}

func (svc *AuthService) ValidateRegistrationToken(ctx context.Context, req *connect.Request[idmv1.ValidateRegistrationTokenRequest]) (*connect.Response[idmv1.ValidateRegistrationTokenResponse], error) {
	_, err := svc.repo.ValidateRegistrationToken(ctx, req.Msg.Token)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.ValidateRegistrationTokenResponse{}), nil
}

func (svc *AuthService) RegisterUser(ctx context.Context, req *connect.Request[idmv1.RegisterUserRequest]) (*connect.Response[idmv1.RegisterUserResponse], error) {
	if !svc.cfg.FeatureEnabled(config.FeatureSelfRegistration) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("registration feature is disabled"))
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userModel, err := svc.CreateUser(ctx, models.User{
		Username: req.Msg.Username,
		Password: string(passwordHash),
	}, req.Msg.RegistrationToken)
	if err != nil {
		return nil, err
	}

	refreshTokenTTL := svc.cfg.RefreshTokenTTL.AsDuration()
	refreshToken, refreshTokenID, err := svc.CreateSignedJWT(*userModel, nil, "", refreshTokenTTL, jwt.ScopeRefresh)
	if err != nil {
		return nil, err
	}

	accessTokenTTL := svc.cfg.AccessTokenTTL.AsDuration()
	accessToken, _, err := svc.CreateSignedJWT(*userModel, nil /*FIXME*/, refreshTokenID, accessTokenTTL, jwt.ScopeAccess)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&idmv1.RegisterUserResponse{
		AccessToken: &idmv1.AccessTokenResponse{
			Token: accessToken,
			User:  conv.UserProtoFromUser(*userModel),
		},
	})

	svc.AddAccessTokenCookie(resp.Header(), accessToken, accessTokenTTL)
	svc.AddRefreshTokenCookie(resp.Header(), refreshToken, refreshTokenTTL)

	return resp, nil
}

func (svc *AuthService) CreateUser(ctx context.Context, userModel models.User, token string) (*models.User, error) {
	var (
		initialRoles []string
		err          error
	)

	// ensure we have a valid registration token if IDM_REGISTRATION_REQUIRES_TOKEN is set to true.
	// Note that we also accept a registration token even if it's not required so users can be
	// bootstrapped with a set of initial roles.
	if svc.cfg.RegistrationRequiresToken || token != "" {
		tokenModel, err := svc.repo.ValidateRegistrationToken(ctx, token)
		if err != nil {
			if errors.Is(err, stmts.ErrNoResults) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid registration token"))
			}
			return nil, err
		}

		if len(tokenModel.InitialRoles) > 0 {
			if err := json.Unmarshal([]byte(tokenModel.InitialRoles), &initialRoles); err != nil {
				return nil, err
			}
		}

		if err := svc.repo.MarkRegistrationTokenUsed(ctx, token); err != nil {
			if errors.Is(err, stmts.ErrNoRowsAffected) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid registration token"))
			}
			return nil, err
		}
	}

	userModel, err = svc.repo.CreateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	merr := new(multierror.Error)
	for _, role := range initialRoles {
		if err := svc.repo.AssignRoleToUser(ctx, userModel.ID, role); err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to assign role %s: %w", role, err))
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return &userModel, err
	}

	return &userModel, nil
}

func (svc *AuthService) AddAccessTokenCookie(resp http.Header, token string, ttl time.Duration) {
	// add the access token as a cookie.
	accessCookie := http.Cookie{
		Name:     svc.cfg.AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		Domain:   svc.cfg.Domain,
		Expires:  time.Now().Add(ttl),
		Secure:   svc.cfg.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	resp.Add("Set-Cookie", accessCookie.String())
}

func (svc *AuthService) AddRefreshTokenCookie(resp http.Header, token string, ttl time.Duration) {
	cookie := http.Cookie{
		Name:     svc.cfg.RefreshTokenCookieName,
		Value:    token,
		Path:     "/tkd.idm.v1.AuthService/RefreshToken",
		Domain:   svc.cfg.Domain,
		Expires:  time.Now().Add(ttl),
		Secure:   svc.cfg.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	resp.Add("Set-Cookie", cookie.String())
}

func (svc *AuthService) CreateSignedJWT(user models.User, roles []models.Role, parentTokenID string, ttl time.Duration, scopes ...jwt.Scope) (string, string, error) {
	auth := &jwt.Authorization{}
	for _, g := range roles {
		auth.Roles = append(auth.Roles, g.ID)
	}

	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(ttl)

	claims := jwt.Claims{
		Audience:  svc.cfg.Audience,
		ExpiresAt: expiresAt.Unix(),
		ID:        tokenID.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    svc.cfg.Domain,
		NotBefore: time.Now().Unix(),
		Subject:   user.ID,
		Name:      user.Username,
		Scopes:    scopes,
		AppMetadata: &jwt.AppMetadata{
			TokenVersion:  "1",
			ParentTokenID: parentTokenID,
			Authorization: auth,
		},
	}

	token, err := jwt.SignToken("HS512", []byte(svc.cfg.JWTSecret), claims)
	if err != nil {
		return "", "", err
	}

	return token, claims.ID, nil
}

func (svc *AuthService) HandleRequestedRedirect(ctx context.Context, requestedRedirect string) (string, error) {
	if requestedRedirect != "" {
		decoded, err := base64.URLEncoding.DecodeString(requestedRedirect)
		if err != nil {
			return "", err
		}

		u, err := url.Parse(string(decoded))
		if err != nil {
			return "", err
		}

		if slices.Contains(svc.cfg.AllowedDomainRedirects, u.Host) {
			middleware.L(ctx).Infof("redirecting user to %s", u.String())
			return u.String(), nil

		} else {
			middleware.L(ctx).Warnf("requested redirect to %s is not allowed", string(decoded))
		}
	}

	return "", nil
}

var _ idmv1connect.AuthServiceHandler = new(AuthService)
