package app

import (
	"net/http"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/apis/pkg/data"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

func (p *Providers) AddRefreshToken(user models.User, roles []models.Role, kind string, headers http.Header) (string, string, error) {
	ttl := p.Config.RefreshTokenTTL.AsDuration()

	for _, overwrite := range p.Config.Overwrites {
		hasUser := slices.Contains(overwrite.UserIDs, user.ID)
		hasRole := data.ElemInBothSlicesFunc(overwrite.RoleIDs, roles, func(r models.Role) string {
			return r.ID
		})

		if (hasUser || hasRole) && overwrite.RefreshTokenTTL.AsDuration() > 0 {
			ttl = overwrite.RefreshTokenTTL.AsDuration()
		}
	}

	signedToken, tokenID, err := p.CreateSignedJWT(user, roles, "", ttl, kind, jwt.ScopeRefresh)
	if err != nil {
		return "", "", err
	}

	if headers != nil {
		p.addRefreshTokenCookie(headers, signedToken)
	}

	return signedToken, tokenID, nil
}

func (p *Providers) AddAccessToken(user models.User, roles []models.Role, ttl time.Duration, parentTokenID string, kind string, headers http.Header) (string, string, error) {
	defaultTTL := p.Config.AccessTokenTTL.AsDuration()

	for _, overwrite := range p.Config.Overwrites {
		hasUser := slices.Contains(overwrite.UserIDs, user.ID)
		hasRole := data.ElemInBothSlicesFunc(overwrite.RoleIDs, roles, func(r models.Role) string {
			return r.ID
		})

		if (hasUser || hasRole) && overwrite.AccessTokenTTL.AsDuration() > 0 {
			defaultTTL = overwrite.AccessTokenTTL.AsDuration()
		}
	}

	if ttl == 0 || ttl > defaultTTL {
		ttl = defaultTTL
	}

	signedToken, tokenID, err := p.CreateSignedJWT(user, roles, parentTokenID, ttl, kind, jwt.ScopeAccess)
	if err != nil {
		return "", "", err
	}

	if headers != nil {
		p.addAccessTokenCookie(headers, signedToken, ttl)
	}

	return signedToken, tokenID, nil

}

func (p *Providers) CreateSignedJWT(user models.User, roles []models.Role, parentTokenID string, ttl time.Duration, kind string, scopes ...jwt.Scope) (string, string, error) {
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
		Audience:    p.Config.Audience,
		ExpiresAt:   expiresAt.Unix(),
		ID:          tokenID.String(),
		IssuedAt:    time.Now().Unix(),
		Issuer:      p.Config.Domain,
		NotBefore:   time.Now().Unix(),
		Subject:     user.ID,
		Name:        user.Username,
		DisplayName: user.DisplayName,
		Scopes:      scopes,
		AppMetadata: &jwt.AppMetadata{
			TokenVersion:  "1",
			ParentTokenID: parentTokenID,
			Authorization: auth,
			LoginKind:     kind,
		},
	}

	token, err := jwt.SignToken("HS512", []byte(p.Config.JWTSecret), claims)
	if err != nil {
		return "", "", err
	}

	return token, claims.ID, nil
}

func (p *Providers) addAccessTokenCookie(resp http.Header, token string, ttl time.Duration) {
	if ttl == 0 {
		ttl = p.Config.AccessTokenTTL.AsDuration()
	}

	// add the access token as a cookie.
	accessCookie := http.Cookie{
		Name:     p.Config.AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		Domain:   p.Config.Domain,
		Expires:  time.Now().Add(ttl),
		Secure:   *p.Config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	resp.Add("Set-Cookie", accessCookie.String())
}

func (p *Providers) addRefreshTokenCookie(resp http.Header, token string) {
	ttl := p.Config.RefreshTokenTTL.AsDuration()

	cookie := http.Cookie{
		Name:     p.Config.RefreshTokenCookieName,
		Value:    token,
		Path:     "/tkd.idm.v1.AuthService/RefreshToken",
		Domain:   p.Config.Domain,
		Expires:  time.Now().Add(ttl),
		Secure:   *p.Config.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	resp.Add("Set-Cookie", cookie.String())
}
