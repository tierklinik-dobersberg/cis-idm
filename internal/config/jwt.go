package config

import (
	"fmt"
	"time"
)

type JWT struct {
	// Audience is the JWT audience that should be used when issuing access tokens.
	Audience string `json:"audience" hcl:"audience,optional"`

	// Secret is the secret that is used to sign access and refresh tokens.
	// Chaning this value during production will invalidate all issued tokens and
	// require all users to re-login.
	Secret string `json:"secret" hcl:"secret"`

	// AccessTokenTTL defines the maximum lifetime for issued access tokens.
	// This defaults to 24h. Users or services requesting an access token
	// may specify a shorter lifetime.
	AccessTokenTTL string `json:"access_token_ttl" hcl:"access_token_ttl,optional"`

	// RefreshTokenTTL defines the lifetime for issued refresh tokens.
	// This defaults to 720h (~1 month)
	RefreshTokenTTL string `json:"refresh_token_ttl" hcl:"refresh_token_ttl,optional"`

	// AccessTokenCookieName is the name of the cookie used to store the
	// access-token for browser requests. This defaults to cis_idm_access.
	AccessTokenCookieName string `json:"access_token_cookie_name" hcl:"access_token_cookie_name,optional"`

	// RefreshTokenCookieName is the name of the cookie used to store the
	// refresh-token for browser requests. This defaults to cis_idm_refresh.
	RefreshTokenCookieName string `json:"refresh_token_cookie_name" hcl:"refresh_token_cookie_name,optional"`

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (file *JWT) ApplyDefaultsAndValidate(domain string) error {
	if file == nil {
		return fmt.Errorf("missing configuration block")
	}

	if file.AccessTokenCookieName == "" {
		file.AccessTokenCookieName = "cis_idm_access"
	}

	if file.RefreshTokenCookieName == "" {
		file.RefreshTokenCookieName = "cis_idm_refresh"
	}

	if file.AccessTokenTTL == "" {
		file.AccessTokenTTL = "1h"
	}

	if ttl, err := time.ParseDuration(file.AccessTokenTTL); err == nil {
		file.accessTokenTTL = ttl
	} else {
		return fmt.Errorf("access_token_ttl: %w", err)
	}

	if file.RefreshTokenTTL == "" {
		file.RefreshTokenTTL = "720h"
	}

	if ttl, err := time.ParseDuration(file.RefreshTokenTTL); err == nil {
		file.refreshTokenTTL = ttl
	} else {
		return fmt.Errorf("refresh_token_ttl: %w", err)
	}

	if file.Secret == "" {
		return fmt.Errorf("missing JWT secret in configuration")
	}

	if file.Audience == "" {
		file.Audience = domain
	}

	return nil
}
