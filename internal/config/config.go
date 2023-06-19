package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Audience               string        `env:"AUDIENCE,required"`
	JWTSecret              string        `env:"JWT_SECRET,required"`
	Domain                 string        `env:"DOMAIN,required"`
	DatabaseURL            string        `env:"RQLITE_URL,required"`
	SecureCookie           bool          `env:"SECURE_COOKIE,default=true"`
	AccessTokenTTL         time.Duration `env:"ACCESS_TOKEN_TTL,default=1m"`
	RefreshTokenTTL        time.Duration `env:"REFRESH_TOKEN_TTL,default=720h"`
	AccessTokenCookieName  string        `env:"ACCESS_TOKEN_COOKIE_NAME,default=cis_idm_access"`
	RefreshTokenCookieName string        `env:"REFRESH_TOKEN_COOKIE_NAME,default=cis_idm_refresh"`
	BootstrapRoles         []string      `env:"BOOTSTRAP_ROLES"`
	LoginRedirectURL       string        `env:"LOGIN_REDIRECT_URL"`
	AllowedDomainRedirects []string      `env:"ALLOWED_DOMAIN_REDIRECTS"`
}

// FromEnvironment returns a Config object parsed from environment variables.
func FromEnvironment(ctx context.Context) (cfg Config, err error) {
	l := envconfig.PrefixLookuper("IDM_", envconfig.OsLookuper())

	if err := envconfig.ProcessWith(ctx, &cfg, l); err != nil {
		return cfg, err
	}

	return cfg, nil
}
