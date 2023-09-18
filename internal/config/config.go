package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/slices"
)

type Twilio struct {
	From        string `json:"from" env:"FROM"`
	AccountSid  string `json:"sid" env:"SID"`
	AccessToken string `json:"token" env:"TOKEN"`
}

type MailConfig struct {
	Host          string `json:"host" env:"HOST"`
	Port          int    `json:"port" env:"PORT"`
	Username      string `json:"user" env:"USER"`
	Password      string `json:"password" env:"PASSWORD"`
	From          string `json:"from" env:"FROM"`
	AllowInsecure bool   `json:"allowInsecure" env:"ALLOW_INSECURE"`
	UseSSL        *bool  `json:"useTLS" env:"USE_TLS"`
}

type Config struct {
	// ForwardAuth configures domains and URLs that require authentication
	// when passed to the /validate endpoint.
	ForwardAuth []*ForwardAuthEntry `json:"forwardAuth" yaml:"forwardAuth"`

	// Audience is the JWT audience that should be used when issuing access tokens.
	Audience string `json:"audience"`

	// JWTSecret is the secret that is used to sign access and refresh tokens.
	// Chaning this value during production will invalidate all issued tokens and
	// require all users to re-login.
	JWTSecret string `json:"jwtSecret"`

	// DatabaseURL is the URL to one of the rqlite cluster members.
	// It should have the format of
	//   http://rqlite:4001/
	DatabaseURL string `json:"rqliteURL"`

	// SecureCookie defines whether or not cookies should be set with the
	// Secure attribute. If left empty, SecureCookie will be automatically
	// set depending on the PublicURL field.
	SecureCookie *bool `json:"secureCookie"`

	// AccessTokenTTL defines the maximum lifetime for issued access tokens.
	// This defaults to 24h. Users or services requesting an access token
	// may specify a shorter lifetime.
	AccessTokenTTL JSONDuration `json:"accessTokenTTL"`

	// RefreshTokenTTL defines the lifetime for issued refresh tokens.
	// This defaults to 720h (~1 month)
	RefreshTokenTTL JSONDuration `json:"refreshTokenTTL"`

	// AccessTokenCookieName is the name of the cookie used to store the
	// access-token for browser requests. This defaults to cis_idm_access.
	AccessTokenCookieName string `json:"accessTokenCookieName"`

	// RefreshTokenCookieName is the name of the cookie used to store the
	// refresh-token for browser requests. This defaults to cis_idm_refresh.
	RefreshTokenCookieName string `json:"refreshTokenCookieName"`

	// BootstrapRoles holds a list of role name that should be automatically
	// created when cisidm is started. Those roles are created with deleteProtection
	// enabled.
	// Use this if you want to ensure cisidm has a set of roles that other services
	// rely upon.
	BootstrapRoles []string `json:"bootstrapRoles"`

	// AllowedDomainRedirects is a list of domain names to which cisidm will allow
	// redirection after login/refresh.
	AllowedDomainRedirects []string `json:"allowedRedirects"`

	// FeatureSet is a list of features that should be enabled. See the AllFeatures
	// global variable for a list of available features. This defaults to "all"
	FeatureSet []Feature `json:"features"`

	// PublicListenAddr defines the listen address for the public listener. This
	// listener requires proper authentication for all endpoints where authentication
	// is specified as required in the protobuf definition.
	// This defaults to :8080
	PublicListenAddr string `json:"publicListener"`

	// AdminListenAddr defines the listen address for the admin listener.
	// All requests received on this listener will automatically get the idm_superuser
	// role assigned. Be careful to not expose this listener to the public!
	// This defaults to :8081
	AdminListenAddr string `json:"adminListener"`

	// AllowedOrigins configures a list of allowed origins for Cross-Origin-Requests.
	// This defaults to the PublicURL as well as http(s)://{{ Domain }}
	AllowedOrigins []string `json:"allowedOrigins"`

	// PublicURL defines the public URL at which cisidm is reachable from the outside.
	// This value MUST be set.
	PublicURL string `json:"publicURL"`

	// StaticFiles defines where cisidm should serve it's user interface from.
	// If left empty, the UI is served from the embedded file-system. If set to
	// a file path than all files from within that directory will be served (see http.Dir
	// for possible security implications). If set to a URL (i.e. starting with "http"),
	// a simple one-host reverse proxy is created.
	// During development, you might want to use `ng serve` from the ui/ folder
	// and set StaticFiles to "http://localhost:4200/"
	StaticFiles string `json:"staticFiles"`

	// RegistrationRequiresToken defines whether or not users are allowed to sign
	// up without a registration token.
	RegistrationRequiresToken bool `json:"registrationRequiresToken"`

	// Domain is the parent domain for which cisidm handles authentication. If you
	// have multiple sub-domains hosting your services you want to set this to the
	// parent domain.
	//
	// I.e. if cisidm is running on account.example.com and you have services on
	// foo.example.com and bar.example.com you want to set the Domain field to "example.com"
	Domain string `json:"domain"`

	// LoginRedirectURL defines the format string to build the redirect URL in the /validate
	// endpoint in case a user needs to authentication.
	// If left empty, it defaults to {{ PublicURL }}/login?redirect=%s
	LoginRedirectURL string `json:"loginURL"`

	// RefreshRedirectURL defines the format string to build the redirect URL in the /validate
	// endpoint in case a user needs to request a new access token.
	// If left empty, it defaults to {{ PublicURL }}/refresh?redirect=%s
	RefreshRedirectURL string `json:"refreshURL"`

	// PasswordResetURL defines the format string to build the password reset URL.
	// If left empty, it defaults to {{ PublicURL }}/password/reset?token=%s
	PasswordResetURL string `json:"passwordResetURL"`

	// VerifyMailURL defines the format string to build the verify-email address URL.
	// If left empty, it defaults to {{ PublicURL }}/verify-mail?token=%s
	VerifyMailURL string `json:"verifyMailURL"`

	// RegistrationURL defines the format string to build the invitation address URL.
	// If left empty, it defaults to {{ PublicURL }}/registration?token=%s
	RegistrationURL string `json:"registrationURL"`

	// SiteName can be used to specify the name of the cisidm instance and will be displayed
	// at the login screen and throughout the user interface. This defaults to Example
	// so will likely want to set this field as well.
	SiteName string `json:"siteName"`

	// SiteNameURL can be set to a URL that will be used to create a HTML link on the login
	// page.
	SiteNameURL string `json:"siteNameUrl"`

	// featureMap is built from FeatureSet and cannot be set using the configuration
	// file.
	featureMap map[Feature]bool `json:"-"`

	// Twilio is required for all SMS related features.
	// TODO(ppacher): print a warning when a SMS feature is enabled
	// but twilio is not confiugred.
	Twilio *Twilio `json:"twilio" envPrefix:"TWILIO__"`

	// MailConfig is required for all email related features.
	MailConfig *MailConfig `json:"mail" envPrefix:"MAIL__"`
}

func LoadFile(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		content, err = yaml.YAMLToJSON(content)
		if err != nil {
			return nil, err
		}

	case ".json":
		// nothing to do here
	default:
		return nil, fmt.Errorf("unsupported file format %q", filepath.Ext(path))
	}

	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()

	var f Config
	if err := dec.Decode(&f); err != nil {
		return nil, err
	}

	if err := f.applyDefaults(); err != nil {
		return &f, nil
	}

	if err := f.parseFeatureSet(); err != nil {
		return &f, err
	}

	return &f, nil
}

func (file *Config) applyDefaults() error {
	if file.PublicURL == "" {
		return fmt.Errorf("publicURL must be set")
	}

	parsedPublicURL, err := url.Parse(file.PublicURL)
	if err != nil {
		return fmt.Errorf("invalid value for publicURL")
	}

	if file.AccessTokenCookieName == "" {
		file.AccessTokenCookieName = "cis_idm_access"
	}
	if file.RefreshTokenCookieName == "" {
		file.RefreshTokenCookieName = "cis_idm_refresh"
	}
	if file.AccessTokenTTL == 0 {
		file.AccessTokenTTL = JSONDuration(time.Hour * 24)
	}
	if file.RefreshTokenTTL == 0 {
		file.RefreshTokenTTL = JSONDuration(time.Hour * 720)
	}
	if file.SiteName == "" {
		file.SiteName = "Example"
	}
	if file.SiteNameURL == "" {
		file.SiteNameURL = "https://example.com"
	}
	if len(file.FeatureSet) == 0 {
		file.FeatureSet = []Feature{FeatureAll}
	}
	if file.PublicListenAddr == "" {
		file.PublicListenAddr = ":8080"
	}
	if file.AdminListenAddr == "" {
		file.AdminListenAddr = ":8081"
	}
	if file.JWTSecret == "" {
		return fmt.Errorf("missing JWT secret in configuration")
	}
	if file.Audience == "" {
		file.Audience = file.Domain
	}

	if file.LoginRedirectURL == "" {
		file.LoginRedirectURL = fmt.Sprintf("%s/login?redirect=%%s", file.PublicURL)
	}

	if file.RefreshRedirectURL == "" {
		file.RefreshRedirectURL = fmt.Sprintf("%s/refresh?redirect=%%s", file.PublicURL)
	}
	if file.PasswordResetURL == "" {
		file.PasswordResetURL = fmt.Sprintf("%s/password/reset?token=%%s", file.PublicURL)
	}
	if file.VerifyMailURL == "" {
		file.VerifyMailURL = fmt.Sprintf("%s/profile/verify-mail?token=%%s", file.PublicURL)
	}
	if file.RegistrationURL == "" {
		file.RegistrationURL = fmt.Sprintf("%s/registration?token=%%s&mail=%%s&name=%%s", file.PublicURL)
	}

	if file.Audience == "" || file.Domain == "" {
		return fmt.Errorf("missing domain and audience")
	}

	if file.SecureCookie == nil {
		file.SecureCookie = new(bool)
		*file.SecureCookie = parsedPublicURL.Scheme == "https:"
	}

	if file.MailConfig == nil {
		file.MailConfig = new(MailConfig)
	}

	return nil
}

func (file Config) AuthRequiredForURL(method string, url string) (*ForwardAuthEntry, bool, error) {
	merr := new(multierror.Error)

	for _, fae := range file.ForwardAuth {
		matches, err := fae.Matches(method, url)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("invalid regex %q: %w", fae.URL, err))
			continue
		}

		if matches {
			if fae.IsRequired() {
				return fae, true, merr.ErrorOrNil()
			}
		}
	}

	return nil, false, merr.ErrorOrNil()
}

// FromEnvironment returns a Config object parsed from environment variables.
func FromEnvironment(ctx context.Context, cfgFilePath string) (cfg Config, err error) {
	parsedFile, err := LoadFile(cfgFilePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file %q: %w", cfgFilePath, err)
	}

	cfg = *parsedFile

	/*
		if err := env.Parse(&cfg); err != nil {
			return cfg, fmt.Errorf("failed to parse config from environment: %w", err)
		}
	*/

	cfg.parseFeatureSet()

	return cfg, nil
}

func (file *Config) parseFeatureSet() error {
	defaultValue := slices.Contains(file.FeatureSet, FeatureAll)

	file.featureMap = make(map[Feature]bool)
	for _, feat := range AllFeatures {
		file.featureMap[feat] = defaultValue
	}

	for _, feat := range file.FeatureSet {
		if feat == FeatureAll {
			continue
		}

		if !slices.Contains(AllFeatures, feat) {
			return fmt.Errorf("unknown feature flag %q", feat)
		}

		allowed := true
		if strings.HasPrefix(string(feat), "!") {
			feat = Feature(strings.TrimPrefix(string(feat), "!"))
			allowed = false
		}

		file.featureMap[feat] = allowed
	}

	return nil
}

func (cfg *Config) FeatureEnabled(feature Feature) bool {
	// if the feature is directly specified return the value
	// directly.
	value, ok := cfg.featureMap[feature]
	if ok {
		return value
	}

	return false
}
