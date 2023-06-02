package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	userdv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/userd/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/userd/v1/userdv1connect"
	"github.com/tierklinik-dobersberg/cis-userd/internal/config"
	"github.com/tierklinik-dobersberg/cis-userd/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-userd/internal/repo/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

type AuthService struct {
	userdv1connect.UnimplementedAuthServiceHandler

	repo UserProvider
	cfg  config.Config
}

type UserProvider interface {
	GetUserByName(ctx context.Context, name string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	GetUserGroupMemberships(ctx context.Context, name string) ([]models.Group, error)
}

// NewService returns a new authentication service that verifies users using repo.
func NewService(repo UserProvider, cfg config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

func (svc *AuthService) Login(ctx context.Context, req *connect.Request[userdv1.LoginRequest]) (*connect.Response[userdv1.LoginResponse], error) {
	logrus.Infof("received authentication request")
	r := req.Msg

	if r.AuthType != userdv1.AuthType_AUTH_TYPE_PASSWORD {
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

	groups, err := svc.repo.GetUserGroupMemberships(ctx, user.ID)
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

	accessToken, err := svc.createAccessToken(user, groups)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := connect.NewResponse(&userdv1.LoginResponse{
		Response: &userdv1.LoginResponse_AccessToken{
			AccessToken: &userdv1.AccessTokenResponse{
				Token: accessToken,
				User: &userdv1.User{
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
			Path:     "/tkd.userd.v1.AuthService/RefreshToken",
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

func (svc *AuthService) Logout(ctx context.Context, req *connect.Request[userdv1.LogoutRequest]) (*connect.Response[userdv1.LogoutResponse], error) {
	return nil, nil
}

func (svc *AuthService) RefreshToken(ctx context.Context, req *connect.Request[userdv1.RefreshTokenRequest]) (*connect.Response[userdv1.RefreshTokenResponse], error) {
	headers := req.Header()
	dummyReq := http.Request{Header: headers}

	var refreshCookie *http.Cookie
	cookies := dummyReq.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == svc.cfg.RefreshTokenCookieName {
			refreshCookie = cookie
			break
		}
	}

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

	groups, err := svc.repo.GetUserGroupMemberships(ctx, claims.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get group memberships: %w", err))
	}

	token, err := svc.createAccessToken(user, groups)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&userdv1.RefreshTokenResponse{
		AccessToken: &userdv1.AccessTokenResponse{
			Token: token,
		},
		// TODO(ppacher): populate .User here?
	}), nil
}

func (svc *AuthService) Introspect(ctx context.Context, req *connect.Request[userdv1.IntrospectRequest]) (*connect.Response[userdv1.IntrospectResponse], error) {

}

func (svc *AuthService) createAccessToken(user models.User, groups []models.Group) (string, error) {
	auth := &jwt.Authorization{}
	for _, g := range groups {
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

var _ userdv1connect.AuthServiceHandler = new(AuthService)
