package config

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slices"
)

type Feature string

const (
	FeatureAll                 = "all"
	FeatureAddresses           = "addresses"
	FeatureEMails              = "emails"
	FeaturePhoneNumbers        = "phoneNumbers"
	FeatureEMailInvite         = "emailInvite"
	FeatureLoginByMail         = "loginByMail"
	FeatureAllowUsernameChange = "allowUsernameChange"
)

var AllFeatures = []Feature{
	FeatureAddresses,
	FeatureEMails,
	FeaturePhoneNumbers,
	FeatureEMailInvite,
	FeatureLoginByMail,
	FeatureAllowUsernameChange,
}

var (
	ErrFeatureDisabled = errors.New("requested feature has been disabled")
)

type Config struct {
	Audience               string        `json:"-" env:"AUDIENCE,required"`
	JWTSecret              string        `json:"-" env:"JWT_SECRET,required"`
	DatabaseURL            string        `json:"-" env:"RQLITE_URL,required"`
	SecureCookie           bool          `json:"-" env:"SECURE_COOKIE,default=true"`
	AccessTokenTTL         time.Duration `json:"-" env:"ACCESS_TOKEN_TTL,default=1m"`
	RefreshTokenTTL        time.Duration `json:"-" env:"REFRESH_TOKEN_TTL,default=720h"`
	AccessTokenCookieName  string        `json:"-" env:"ACCESS_TOKEN_COOKIE_NAME,default=cis_idm_access"`
	RefreshTokenCookieName string        `json:"-" env:"REFRESH_TOKEN_COOKIE_NAME,default=cis_idm_refresh"`
	BootstrapRoles         []string      `json:"-" env:"BOOTSTRAP_ROLES"`
	AllowedDomainRedirects []string      `json:"-" env:"ALLOWED_DOMAIN_REDIRECTS"`

	// Exposed via /config endpoint
	Domain           string    `json:"domain" env:"DOMAIN,required"`
	LoginRedirectURL string    `json:"loginURL" env:"LOGIN_REDIRECT_URL"`
	SiteName         string    `json:"siteName" env:"SITE_NAME,default=Example"`
	SiteNameURL      string    `json:"siteNameUrl" env:"SITE_NAME_URL"`
	FeatureSet       []Feature `json:"-" env:"ENABLED_FEATURES,default=all"`

	FeatureMap map[Feature]bool `json:"features"`
}

// FromEnvironment returns a Config object parsed from environment variables.
func FromEnvironment(ctx context.Context) (cfg Config, err error) {
	l := envconfig.PrefixLookuper("IDM_", envconfig.OsLookuper())

	if err := envconfig.ProcessWith(ctx, &cfg, l); err != nil {
		return cfg, err
	}

	defaultValue := slices.Contains(cfg.FeatureSet, FeatureAll)

	cfg.FeatureMap = make(map[Feature]bool)
	for _, feat := range AllFeatures {
		cfg.FeatureMap[feat] = defaultValue
	}

	for _, feat := range cfg.FeatureSet {
		allowed := true
		if strings.HasPrefix(string(feat), "!") {
			feat = Feature(strings.TrimPrefix(string(feat), "!"))
			allowed = false
		}

		cfg.FeatureMap[feat] = allowed
	}

	return cfg, nil
}

func (cfg *Config) FeatureEnabled(feature Feature) bool {
	// if the feature is directly specified return the value
	// directly.
	value, ok := cfg.FeatureMap[feature]
	if ok {
		return value
	}

	return false
}
