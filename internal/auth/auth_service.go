package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type AuthService struct {
	idmv1connect.UnimplementedAuthServiceHandler

	validator *protovalidate.Validator

	repo UserProvider
	cfg  config.Config
}

type UserProvider interface {
	GetUserByName(ctx context.Context, name string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	GetUserRoles(ctx context.Context, name string) ([]models.Role, error)
	MarkTokenRejected(ctx context.Context, token models.RejectedToken) error
}

// NewService returns a new authentication service that verifies users using repo.
func NewService(repo UserProvider, cfg config.Config) (*AuthService, error) {
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

	if r.AuthType != idmv1.AuthType_AUTH_TYPE_PASSWORD {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("auth type not supported"))
	}

	passwordAuth := r.GetPassword()
	if passwordAuth == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid payload for password auth type"))
	}

	logrus.Infof("authentication request for user %s", passwordAuth.GetUsername())

	user, err := svc.repo.GetUserByName(ctx, passwordAuth.GetUsername())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not found: %w", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordAuth.GetPassword())); err != nil {
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("incorrect password"))
		}
	}

	roles, err := svc.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

	accessToken, err := svc.createAccessToken(user, roles)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
	})

	if !r.GetNoRefreshToken() {
		refreshToken, err := svc.createRefreshToken(user)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		cookie := http.Cookie{
			Name:     svc.cfg.RefreshTokenCookieName,
			Value:    refreshToken,
			Path:     "/tkd.idm.v1.AuthService/RefreshToken",
			Domain:   svc.cfg.Domain,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		resp.Header().Add("Set-Cookie", cookie.String())
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
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no claims associated with request context"))
	}

	// mark the token as rejected
	if err := svc.repo.MarkTokenRejected(ctx, models.RejectedToken{
		TokenID:  claims.ID,
		UserID:   claims.Subject,
		IssuedAt: claims.IssuedAt,
		ExiresAt: claims.ExpiresAt,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to mark token as rejected: %w", err))
	}

	return connect.NewResponse(new(idmv1.LogoutResponse)), nil
}

func (svc *AuthService) RefreshToken(ctx context.Context, req *connect.Request[idmv1.RefreshTokenRequest]) (*connect.Response[idmv1.RefreshTokenResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	refreshCookie := getRefreshTokenCookie(svc.cfg.RefreshTokenCookieName, req.Header())
	if refreshCookie == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no refresh cookie provided"))
	}

	if err := refreshCookie.Valid(); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	claims, err := jwt.ParseAndVerify([]byte(svc.cfg.JWTSecret), refreshCookie.Value)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
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
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get group memberships: %w", err))
	}

	token, err := svc.createAccessToken(user, roles)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&idmv1.RefreshTokenResponse{
		AccessToken: &idmv1.AccessTokenResponse{
			Token: token,
			User:  conv.UserProtoFromUser(user),
		},
	}), nil
}

func (svc *AuthService) Introspect(ctx context.Context, req *connect.Request[idmv1.IntrospectRequest]) (*connect.Response[idmv1.IntrospectResponse], error) {
	if err := svc.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("not jwt token claims found in request context"))
	}

	user, err := svc.repo.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("invalid user"))
	}

	return connect.NewResponse(&idmv1.IntrospectResponse{
		User: conv.UserProtoFromUser(user),
	}), nil
}

func (svc *AuthService) createAccessToken(user models.User, roles []models.Role) (string, error) {
	auth := &jwt.Authorization{}
	for _, g := range roles {
		auth.Roles = append(auth.Roles, g.ID)
	}

	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(svc.cfg.AccessTokenTTL)

	claims := jwt.Claims{
		Audience:  "dobersberg.vet",
		ExpiresAt: expiresAt.Unix(),
		ID:        tokenID.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    svc.cfg.Domain,
		NotBefore: time.Now().Unix(),
		Subject:   user.ID,
		Name:      user.Username,
		Scopes:    []jwt.Scope{jwt.ScopeAccess},
		AppMetadata: &jwt.AppMetadata{
			TokenVersion:  "1",
			Authorization: auth,
		},
	}

	token, err := jwt.SignToken("HS512", []byte(svc.cfg.JWTSecret), claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *AuthService) createRefreshToken(user models.User) (string, error) {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(svc.cfg.RefreshTokenTTL)

	claims := jwt.Claims{
		Audience:  "dobersberg.vet",
		ExpiresAt: expiresAt.Unix(),
		ID:        tokenID.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    svc.cfg.Domain,
		NotBefore: time.Now().Unix(),
		Subject:   user.ID,
		Name:      user.Username,
		Scopes:    []jwt.Scope{jwt.ScopeRefresh},
	}

	token, err := jwt.SignToken("HS512", []byte(svc.cfg.JWTSecret), claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func getRefreshTokenCookie(cookieName string, headers http.Header) *http.Cookie {
	// we create a dummy http request so we can use the cookie parser
	// from the stdlib which is, unfortunately, not exported for direct
	// use.
	dummyReq := http.Request{Header: headers}

	cookies := dummyReq.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie
		}
	}

	return nil
}

var _ idmv1connect.AuthServiceHandler = new(AuthService)
