// Inspired by and designed to work with
// https://github.com/greenpau/caddy-auth-jwt/

// Package jwt provides JWT token signing.
package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Scope defines the scope of a JWT token.
type Scope string

const (
	// ScopeAccess is a full access token.
	ScopeAccess = "access"

	// ScopeRefresh is required to receive a new access token.
	ScopeRefresh = "refresh"

	// Scope2FAPending is used for JWTs that are issued during the login
	// process when the second authentication factor is still pending.
	Scope2FAPending = "2fa-pending"
)

var supportedMethods = map[string]struct{}{
	"HS512": {},
	"HS384": {},
	"HS256": {},
}

// Authorization contains app related authorization and permission
// settings.
type Authorization struct {
	Roles []string `json:"roles,omitempty" xml:"roles" yaml:"roles,omitempty"`
}

type LoginKind string

const (
	LoginKindInvalid  LoginKind = ""
	LoginKindPassword LoginKind = "password"
	LoginKindMFA      LoginKind = "mfa"
	LoginKindWebauthn LoginKind = "webauthn"
	LoginKindAPI      LoginKind = "api"
)

// AppMetadata defines app specific metadata attached to
// JWT tokens issued by cisd.
type AppMetadata struct {
	TokenVersion  string         `json:"token_version" xml:"token_version" yaml:"token_version"`
	ParentTokenID string         `json:"parent_token" xml:"parent_token" yaml:"parent_token"`
	Authorization *Authorization `json:"authorization,omitempty" xml:"authorization" yaml:"authorization,omitempty"`
	LoginKind     LoginKind      `json:"loginKind,omitempty"`
}

// Claims represents the claims added to a JWT token issued
// by cisd.
type Claims struct {
	Audience    string       `json:"aud,omitempty" xml:"aud" yaml:"aud,omitempty"`
	ExpiresAt   int64        `json:"exp,omitempty" xml:"exp" yaml:"exp,omitempty"`
	ID          string       `json:"jti,omitempty" xml:"jti" yaml:"jti,omitempty"`
	IssuedAt    int64        `json:"iat,omitempty" xml:"iat" yaml:"iat,omitempty"`
	Issuer      string       `json:"iss,omitempty" xml:"iss" yaml:"iss,omitempty"`
	NotBefore   int64        `json:"nbf,omitempty" xml:"nbf" yaml:"nbf,omitempty"`
	Subject     string       `json:"sub,omitempty" xml:"sub" yaml:"sub,omitempty"`
	Name        string       `json:"name,omitempty" xml:"name" yaml:"name,omitempty"`
	DisplayName string       `json:"displayName,omitempty" xml:"displayName" yaml:"displayName"`
	Scopes      []Scope      `json:"scopes,omitempty" xml:"scopes" yaml:"scopes,omitempty"`
	Email       string       `json:"email,omitempty" xml:"email" yaml:"email,omitempty"`
	AppMetadata *AppMetadata `json:"app_metadata,omitempty" xml:"app_metadata" yaml:"app_metadata,omitempty"`
}

// Valid returns true if the token is valid and can be used.
// It checks if the token has been expired or may not yet be
// used.
func (u Claims) Valid() error {
	if u.NotBefore > 0 && time.Now().Unix() < u.NotBefore {
		return jwt.NewValidationError("Not yet valid", jwt.ValidationErrorNotValidYet)
	}

	if u.ExpiresAt > 0 && u.ExpiresAt < time.Now().Unix() {
		return jwt.NewValidationError("Expired", jwt.ValidationErrorExpired)
	}

	return nil
}

// Sign returns a signed JWT token from u.
func (u *Claims) Sign(method string, secret []byte) (string, error) {
	return SignToken(method, secret, *u)
}

// SignToken returns a signed JWT token.
func SignToken(method string, secret []byte, claims Claims) (string, error) {
	if _, exists := supportedMethods[method]; !exists {
		return "", fmt.Errorf("unsupported signing method %q", method)
	}

	if secret == nil {
		return "", fmt.Errorf("missing secret")
	}

	sm := jwt.GetSigningMethod(method)
	token := jwt.NewWithClaims(sm, claims)
	signedToken, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseAndVerify parses the JWT token and verifies it's signature.
func ParseAndVerify(secret []byte, token string) (*Claims, error) {
	var c Claims

	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return &c, err
	}

	return &c, nil
}
