package config

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/permission"
	"golang.org/x/exp/slices"
)

type Twilio struct {
	From        string `json:"from" hcl:"from"`
	AccountSid  string `json:"sid" hcl:"sid"`
	AccessToken string `json:"token" hcl:"token"`
}

type MailConfig struct {
	Host          string `json:"host" hcl:"host"`
	Port          int    `json:"port" hcl:"port"`
	Username      string `json:"user" hcl:"user"`
	Password      string `json:"password" hcl:"password"`
	From          string `json:"from" hcl:"from"`
	AllowInsecure bool   `json:"allow_insecure" hcl:"allow_insecure,optional"`
	UseSSL        *bool  `json:"use_tls" hcl:"use_tls,optional"`
}

type DryRun struct {
	MailTarget string `json:"mail" hcl:"mail,optional"`
	SMSTarget  string `json:"sms" hcl:"sms,optional"`
}

type Overwrite struct {
	UserIDs []string `json:"users" hcl:"user_ids,optional"`
	RoleIDs []string `json:"roles" hcl:"role_ids,optional"`

	AccessTokenTTL  JSONDuration `json:"access_token_ttl" hcl:"access_token_ttl,optional"`
	RefreshTokenTTL JSONDuration `json:"refresh_token_ttl" hcl:"refresh_token_ttl,optional"`
}

type WebPush struct {
	Admin           string `json:"admin" hcl:"admin"`
	VAPIDpublicKey  string `json:"vapid_public_key" hcl:"vapid_public_key"`
	VAPIDprivateKey string `json:"vapid_private_key" hcl:"vapid_private_key"`
}

type Role struct {
	ID          string   `json:"id" hcl:",label"`
	Name        string   `json:"name" hcl:"name"`
	Description string   `json:"description" hcl:"description,optional"`
	Permissions []string `json:"permissions" hcl:"permissions,optional"`
}

type Config struct {
	// LogLevel defines the log level to use.
	LogLevel string `json:"log_level" hcl:"log_level,optional"`

	// ForwardAuth configures domains and URLs that require authentication
	// when passed to the /validate endpoint.
	ForwardAuth []*ForwardAuthEntry `json:"forward_auth" hcl:"forward_auth,block"`

	// DryRun may be set to enable dry-run mode which allows overwriting
	// notification targets.
	DryRun *DryRun `json:"dry_run" hcl:"dry_run,block"`

	// TrustedNetworks is a list of CIDR network addresses that are considered
	// trusted. Any X-Forwareded-For header from these networks will be parsed
	// and applied.
	TrustedNetworks []string `json:"trusted_networks" hcl:"trusted_networks,optional"`

	// Audience is the JWT audience that should be used when issuing access tokens.
	Audience string `json:"audience" hcl:"audience,optional"`

	// JWTSecret is the secret that is used to sign access and refresh tokens.
	// Chaning this value during production will invalidate all issued tokens and
	// require all users to re-login.
	JWTSecret string `json:"jwt_secret" hcl:"jwt_secret"`

	DatabaseURL string `json:"database_url" hcl:"database_url"`

	// SecureCookie defines whether or not cookies should be set with the
	// Secure attribute. If left empty, SecureCookie will be automatically
	// set depending on the PublicURL field.
	SecureCookie *bool `json:"secure_cookies" hcl:"secure_cookies,optional"`

	// AccessTokenTTL defines the maximum lifetime for issued access tokens.
	// This defaults to 24h. Users or services requesting an access token
	// may specify a shorter lifetime.
	AccessTokenTTL JSONDuration `json:"access_token_ttl" hcl:"access_token_ttl,optional"`

	// RefreshTokenTTL defines the lifetime for issued refresh tokens.
	// This defaults to 720h (~1 month)
	RefreshTokenTTL JSONDuration `json:"refresh_token_ttl" hcl:"refresh_token_ttl,optional"`

	// AccessTokenCookieName is the name of the cookie used to store the
	// access-token for browser requests. This defaults to cis_idm_access.
	AccessTokenCookieName string `json:"access_token_cookie_name" hcl:"access_token_cookie_name,optional"`

	// RefreshTokenCookieName is the name of the cookie used to store the
	// refresh-token for browser requests. This defaults to cis_idm_refresh.
	RefreshTokenCookieName string `json:"refresh_token_cookie_name" hcl:"refresh_token_cookie_name,optional"`

	// Overwrites may hold configuration overwrites per user or role.
	Overwrites []Overwrite `json:"overwrites" hcl:"overwrites,block"`

	// Roles holds a list of role name that should be automatically
	// created when cisidm is started. Those roles are created with deleteProtection
	// enabled.
	// Use this if you want to ensure cisidm has a set of roles that other services
	// rely upon.
	Roles []Role `json:"roles" hcl:"role,block"`

	// EnableDynamicRoles controles whether or not roles can be created/updated/deleted
	// via the tkd.idm.v1.RoleService API.
	// This defaults to true if no roles are configured in the configuration file, otherwise,
	// if roles are pre-configured, this defaults to false.
	// To have config defined roles while still allowing role management via the API you need
	// to explicitly set EnableDynamicRoles to true.
	//
	// Note that even if this is enabled, roles configured via the configuration file cannot
	// be modified or deleted.
	EnableDynamicRoles *bool `json:"enable_dynamic_roles" hcl:"enable_dynamic_roles,optional"`

	// Permissions defines the hierarchical set of available permissions.
	// Note that the specified permission tree will be merged into the default set of permissions
	// that are built into cisidm.
	Permissions []string `json:"permissions" hcl:"permissions,optional"`

	// AllowedDomainRedirects is a list of domain names to which cisidm will allow
	// redirection after login/refresh.
	AllowedDomainRedirects []string `json:"allowed_redirects" hcl:"allowed_redirects,optional"`

	// FeatureSet is a list of features that should be enabled. See the AllFeatures
	// global variable for a list of available features. This defaults to "all"
	FeatureSet []Feature `json:"features" hcl:"features,optional"`

	// PublicListenAddr defines the listen address for the public listener. This
	// listener requires proper authentication for all endpoints where authentication
	// is specified as required in the protobuf definition.
	// This defaults to :8080
	PublicListenAddr string `json:"public_listener" hcl:"public_listener,optional"`

	// AdminListenAddr defines the listen address for the admin listener.
	// All requests received on this listener will automatically get the idm_superuser
	// role assigned. Be careful to not expose this listener to the public!
	// This defaults to :8081
	AdminListenAddr string `json:"admin_listener" hcl:"admin_listener,optional"`

	// AllowedOrigins configures a list of allowed origins for Cross-Origin-Requests.
	// This defaults to the PublicURL as well as http(s)://{{ Domain }}
	AllowedOrigins []string `json:"allowed_origins" hcl:"allowed_origins,optional"`

	// PublicURL defines the public URL at which cisidm is reachable from the outside.
	// This value MUST be set.
	PublicURL string `json:"public_url" hcl:"public_url"`

	// StaticFiles defines where cisidm should serve it's user interface from.
	// If left empty, the UI is served from the embedded file-system. If set to
	// a file path than all files from within that directory will be served (see http.Dir
	// for possible security implications). If set to a URL (i.e. starting with "http"),
	// a simple one-host reverse proxy is created.
	// During development, you might want to use `ng serve` from the ui/ folder
	// and set StaticFiles to "http://localhost:4200/"
	StaticFiles string `json:"static_files" hcl:"static_files,optional"`

	// ExtraAssetsDirectory can be set to a directory (or HTTP URL)
	// that will be used to serve additional files at the /files endpoint.
	ExtraAssetsDirectory string `json:"extra_assets" hcl:"extra_assets,optional"`

	// LogoURL may be set to a path or HTTP resource that should be displayed as the
	// application logo on the login screen.
	LogoURL string `json:"logo_url" hcl:"logo_url,optional"`

	// RegistrationRequiresToken defines whether or not users are allowed to sign
	// up without a registration token.
	RegistrationRequiresToken bool `json:"registration_requires_token" hcl:"registration_requires_token,optional"`

	// Domain is the parent domain for which cisidm handles authentication. If you
	// have multiple sub-domains hosting your services you want to set this to the
	// parent domain.
	//
	// I.e. if cisidm is running on account.example.com and you have services on
	// foo.example.com and bar.example.com you want to set the Domain field to "example.com"
	Domain string `json:"domain" hcl:"domain"`

	// LoginRedirectURL defines the format string to build the redirect URL in the /validate
	// endpoint in case a user needs to authentication.
	// If left empty, it defaults to {{ PublicURL }}/login?redirect=%s
	LoginRedirectURL string `json:"login_url" hcl:"login_url,optional"`

	// RefreshRedirectURL defines the format string to build the redirect URL in the /validate
	// endpoint in case a user needs to request a new access token.
	// If left empty, it defaults to {{ PublicURL }}/refresh?redirect=%s
	RefreshRedirectURL string `json:"refresh_url" hcl:"refresh_url,optional"`

	// PasswordResetURL defines the format string to build the password reset URL.
	// If left empty, it defaults to {{ PublicURL }}/password/reset?token=%s
	PasswordResetURL string `json:"password_reset_url" hcl:"password_reset_url,optional"`

	// VerifyMailURL defines the format string to build the verify-email address URL.
	// If left empty, it defaults to {{ PublicURL }}/verify-mail?token=%s
	VerifyMailURL string `json:"verify_mail_url" hcl:"verify_mail_url,optional"`

	// RegistrationURL defines the format string to build the invitation address URL.
	// If left empty, it defaults to {{ PublicURL }}/registration?token=%s
	RegistrationURL string `json:"registration_url" hcl:"registration_url,optional"`

	// SiteName can be used to specify the name of the cisidm instance and will be displayed
	// at the login screen and throughout the user interface. This defaults to Example
	// so will likely want to set this field as well.
	SiteName string `json:"site_name" hcl:"site_name"`

	// SiteNameURL can be set to a URL that will be used to create a HTML link on the login
	// page.
	SiteNameURL string `json:"site_name_url" hcl:"site_name_url,optional"`

	// Twilio is required for all SMS related features.
	// TODO(ppacher): print a warning when a SMS feature is enabled
	// but twilio is not confiugred.
	Twilio *Twilio `json:"twilio" hcl:"twilio,block"`

	// MailConfig is required for all email related features.
	MailConfig *MailConfig `json:"mail" hcl:"mail,block"`

	// ExtraDataConfig defines the schema and visibility for the user extra data.
	ExtraDataConfig []*FieldConfig `json:"extra_data" hcl:"extra_data,block"`

	// WebPush holds VAPID keys for web-push integration.
	WebPush *WebPush `json:"webpush" hcl:"webpush,block"`

	// featureMap is built from FeatureSet and cannot be set using the configuration
	// file.
	featureMap map[Feature]bool `json:"-"`

	permissionTree permission.Tree
}

func LoadFile(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".yml", ".yaml":
		content, err = yaml.YAMLToJSON(content)
		if err != nil {
			return nil, err
		}

		ext = ".json"

	case ".json", ".hcl":
		// nothing to do here
	default:
		return nil, fmt.Errorf("unsupported file format %q", filepath.Ext(path))
	}

	/*
		dec := json.NewDecoder(bytes.NewReader(content))
		dec.DisallowUnknownFields()

		var f Config
		if err := dec.Decode(&f); err != nil {
			return nil, err
		}
	*/

	var f Config
	var ctx hcl.EvalContext

	if err := hclsimple.Decode(filepath.Base(path)+ext, content, &ctx, &f); err != nil {
		return &f, err
	}

	if err := f.applyDefaults(); err != nil {
		return &f, nil
	}

	if err := f.parseFeatureSet(); err != nil {
		return &f, err
	}

	f.permissionTree = permission.Tree{}
	for _, p := range f.Permissions {
		f.permissionTree.Insert(p)
	}

	return &f, nil
}

func (file *Config) PermissionTree() permission.Tree {
	return file.permissionTree
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

	if file.EnableDynamicRoles == nil {
		b := len(file.Roles) == 0
		file.EnableDynamicRoles = &b
	}

	// validate the user extra data.
	if len(file.ExtraDataConfig) > 0 {
		for key, cfg := range file.ExtraDataConfig {
			if err := cfg.ValidateConfig(FieldVisibilityPublic); err != nil {
				return fmt.Errorf("extraData: %s: %w", key, err)
			}
		}
	}

	return nil
}

func (file Config) DynmicRolesEnabled() bool {
	if file.EnableDynamicRoles == nil {
		return false
	}

	return *file.EnableDynamicRoles
}

func (file Config) AuthRequiredForURL(ctx context.Context, method string, url string) (*ForwardAuthEntry, bool, error) {
	merr := new(multierror.Error)
	l := log.L(ctx)

	for idx, fae := range file.ForwardAuth {
		matches, err := fae.Matches(method, url)
		if err != nil {
			l.Debugf("forward-auth[%d] failed to match: %s", idx, err)

			merr.Errors = append(merr.Errors, fmt.Errorf("invalid regex %q: %w", fae.URL, err))

			continue
		}

		if matches {
			l.Debugf("forward-auth[%d] entry matches request to %s %s", idx, method, url)

			return fae, fae.IsRequired(), merr.ErrorOrNil()
		} else {
			l.Debugf("forward-auth[%d] entry does not match request to %s %s", idx, method, url)
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

		allowed := true
		if strings.HasPrefix(string(feat), "!") {
			feat = Feature(strings.TrimPrefix(string(feat), "!"))
			allowed = false
		}

		if !slices.Contains(AllFeatures, feat) {
			return fmt.Errorf("unknown feature flag %q", feat)
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
