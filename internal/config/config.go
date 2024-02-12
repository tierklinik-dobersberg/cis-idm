package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
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
	Type string `json:"type" hcl:"type,label"` // role or user
	ID   string `json:"id" hcl:"id,label"`

	AccessTokenTTL  string `json:"access_token_ttl" hcl:"access_token_ttl,optional"`
	RefreshTokenTTL string `json:"refresh_token_ttl" hcl:"refresh_token_ttl,optional"`

	accessTTL  time.Duration
	refreshTTL time.Duration
}

func (ov *Overwrite) Validate() error {
	var err error

	ov.accessTTL, err = time.ParseDuration(ov.AccessTokenTTL)
	if err != nil {
		return fmt.Errorf("access_token_ttl: %w", err)
	}

	ov.refreshTTL, err = time.ParseDuration(ov.RefreshTokenTTL)
	if err != nil {
		return fmt.Errorf("refresh_token_ttl: %w", err)
	}

	return nil
}

func (ov *Overwrite) AccessTTL() time.Duration  { return ov.accessTTL }
func (ov *Overwrite) RefreshTTL() time.Duration { return ov.refreshTTL }

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

type RegistrationMode string

const (
	RegistrationModePublic   = RegistrationMode("public")
	RegistrationModeToken    = RegistrationMode("token")
	RegistrationModeDisabled = RegistrationMode("disabled")
)

type Policy struct {
	Name    string `json:"name" hcl:"name,label"`
	Content string `json:"content" hcl:"content"`
}

type PolicyConfig struct {
	Directories []string `json:"directories" hcl:"directories,optional"`
	Debug       bool     `json:"debug" hcl:"debug,optional"`

	Policies []Policy `json:"policy" hcl:"policy,block"`
}

type ForwardAuthConfig struct {
	// RegoQuery is the rego policy query that cis-idm should perform
	// when evaluating forward auth policies.
	// Defaults to "data.cisidm.forward_auth"
	RegoQuery string `json:"rego_query" hcl:"rego_query,optional"`

	// Default holds the default policy for forward auth queries.
	// This may either be set to "allow" or "deny" (default).
	//
	// Depending on the value of Default cisidm will look for different rules
	// when evaluating policies.
	// If Default is set to "allow", cisidm will evaluate any "deny" rule.
	// If Default is set to "deny", cisidm will evaluate any "allow" rule.
	Default string `json:"default" hcl:"default,optional"`

	// AllowCORSPreflight might be set to enable or disable automatic pass-through of
	// CORS preflight requests.
	// Defaults to true.
	AllowCORSPreflight *bool `json:"allow_cors_preflight" hcl:"allow_cors_preflight,optional"`

	UserIDHeader             *string `json:"user_id_header" hcl:"user_id_header,optional"`
	UsernameHeader           *string `json:"username_header" hcl:"username_header,optional"`
	MailHeader               *string `json:"mail_header" hcl:"mail_header,optional"`
	RoleHeader               *string `json:"role_header" hcl:"role_header,optional"`
	AvatarURLHeader          *string `json:"avatar_url_header" hcl:"avatar_url_header,optional"`
	DisplayNameHeader        *string `json:"display_name_header" hcl:"display_name_header,optional"`
	ResolvedPermissionHeader *string `json:"permission_header" hcl:"permission_header,optional"`
}

var (
	defaultUserIDHeader      = "X-Remote-User-ID"
	defaultUsernameHeader    = "X-Remote-User"
	defaultMailHeader        = "X-Remote-Mail"
	defaultRoleHeader        = "X-Remote-Role"
	defaultDisplayNameHeader = "X-Remote-User-Display-Name"
	defaultAvatarURLHeader   = "X-Remote-Avatar-URL"
	defaultPermissionHeader  = "X-Remote-Permission"
)

func (cfg *PolicyConfig) ApplyDefaultsAndValidate() error {
	if cfg == nil {
		return nil
	}

	return nil
}

func (cfg *ForwardAuthConfig) ApplyDefaultsAndValidate() error {
	if cfg == nil {
		return nil
	}

	switch cfg.Default {
	case "":
		cfg.Default = "deny"
	case "allow", "deny":
	default:
		return fmt.Errorf("default: invalid value, expected \"allow\" or \"deny\"")
	}

	if cfg.RegoQuery == "" {
		cfg.RegoQuery = "data.cisidm.forward_auth"
	}

	if cfg.AllowCORSPreflight == nil {
		val := true
		cfg.AllowCORSPreflight = &val
	}

	if cfg.UserIDHeader == nil {
		cfg.UserIDHeader = &defaultUserIDHeader
	}

	if cfg.UsernameHeader == nil {
		cfg.UsernameHeader = &defaultUsernameHeader
	}

	if cfg.MailHeader == nil {
		cfg.MailHeader = &defaultMailHeader
	}

	if cfg.RoleHeader == nil {
		cfg.RoleHeader = &defaultRoleHeader
	}

	if cfg.AvatarURLHeader == nil {
		cfg.AvatarURLHeader = &defaultAvatarURLHeader
	}

	if cfg.DisplayNameHeader == nil {
		cfg.DisplayNameHeader = &defaultDisplayNameHeader
	}

	if cfg.ResolvedPermissionHeader == nil {
		cfg.ResolvedPermissionHeader = &defaultPermissionHeader
	}

	return nil
}

type Config struct {
	// LogLevel defines the log level to use.
	LogLevel string `json:"log_level" hcl:"log_level,optional"`

	// PolicyConfig holds the configuration for rego policies.
	PolicyConfig PolicyConfig `json:"policy" hcl:"policies,block"`

	// ForwardAuth provides forward auth configuration.
	ForwardAuth ForwardAuthConfig `json:"forward_auth" hcl:"forward_auth,block"`

	// Server holds the server configuration block include CORS, listen addresses
	// and cookie settings.
	Server *Server `hcl:"server,block" json:"server"`

	// JWT holds the JWT configuration.
	JWT *JWT `hcl:"jwt,block" json:"jwt"`

	// UserInterface configures settings for user facing interfaces like the built-in
	// Web-Interface or mail/SMS templates.
	UserInterface *UserInterface `hcl:"ui,block" json:"ui"`

	// DryRun may be set to enable dry-run mode which allows overwriting
	// notification targets.
	DryRun *DryRun `json:"dry_run" hcl:"dry_run,block"`

	DatabaseURL string `json:"database_url" hcl:"database_url"`

	// Roles holds a list of role name that should be automatically
	// created when cisidm is started. Those roles are created with deleteProtection
	// enabled.
	// Use this if you want to ensure cisidm has a set of roles that other services
	// rely upon.
	Roles []Role `json:"roles" hcl:"role,block"`

	// Overwrites may hold configuration overwrites per user or role.
	Overwrites []Overwrite `json:"overwrite" hcl:"overwrite,block"`

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

	// PermissionTrees may be set to true to enable permission trees.
	PermissionTrees bool `json:"permission_trees" hcl:"permission_trees,optional"`

	// Permissions defines the hierarchical set of available permissions.
	// Note that the specified permission tree will be merged into the default set of permissions
	// that are built into cisidm.
	Permissions []string `json:"permissions" hcl:"permissions,optional"`

	// FeatureSet is a list of features that should be enabled. See the AllFeatures
	// global variable for a list of available features. This defaults to "all"
	FeatureSet []Feature `json:"features" hcl:"features,optional"`

	// RegistrationMode defines whether or not users are allowed to sign
	// up without a registration token.
	RegistrationMode RegistrationMode `json:"registration" hcl:"registration,optional"`

	// Twilio is required for all SMS related features.
	// TODO(ppacher): print a warning when a SMS feature is enabled
	// but twilio is not confiugred.
	Twilio *Twilio `json:"twilio" hcl:"twilio,block"`

	// MailConfig is required for all email related features.
	MailConfig *MailConfig `json:"mail" hcl:"mail,block"`

	// ExtraDataConfig defines the schema and visibility for the user extra data.
	ExtraDataConfig []*FieldConfig `json:"field" hcl:"field,block"`

	// WebPush holds VAPID keys for web-push integration.
	WebPush *WebPush `json:"webpush" hcl:"webpush,block"`

	// featureMap is built from FeatureSet and cannot be set using the configuration
	// file.
	featureMap map[Feature]bool `json:"-"`

	permissionTree permission.Resolver
}

func (cfg *Config) AccessTTL() time.Duration {
	return cfg.JWT.accessTokenTTL
}

func (cfg *Config) RefreshTTL() time.Duration {
	return cfg.JWT.refreshTokenTTL
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

	if f.PermissionTrees {
		tree := permission.Tree{}
		for _, p := range f.Permissions {
			tree.Insert(p)
		}

		f.permissionTree = tree
	} else {
		f.permissionTree = permission.NoTree{}
	}

	return &f, nil
}

func (file *Config) PermissionTree() permission.Resolver {
	return file.permissionTree
}

func (file *Config) applyDefaults() error {
	if err := file.UserInterface.ApplyDefaultsAndValidate(); err != nil {
		return fmt.Errorf("ui: %w", err)
	}

	parsedPublicURL, err := url.Parse(file.UserInterface.PublicURL)
	if err != nil {
		return fmt.Errorf("ui: invalid value for publicURL")
	}

	if err := file.Server.ApplyDefaultsAndValidate(parsedPublicURL.Scheme == "https:"); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	if err := file.JWT.ApplyDefaultsAndValidate(file.Server.Domain); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}

	if err := file.PolicyConfig.ApplyDefaultsAndValidate(); err != nil {
		return fmt.Errorf("policies: %w", err)
	}

	if err := file.ForwardAuth.ApplyDefaultsAndValidate(); err != nil {
		return fmt.Errorf("forward_auth: %w", err)
	}

	if len(file.FeatureSet) == 0 {
		file.FeatureSet = []Feature{FeatureAll}
	}

	if file.MailConfig == nil {
		file.MailConfig = new(MailConfig)
	}

	if file.EnableDynamicRoles == nil {
		b := len(file.Roles) == 0
		file.EnableDynamicRoles = &b
	}

	if file.RegistrationMode == "" {
		file.RegistrationMode = RegistrationModeDisabled
	}

	// validate the user extra data.
	if len(file.ExtraDataConfig) > 0 {
		for key, cfg := range file.ExtraDataConfig {
			if err := cfg.ValidateConfig(FieldVisibilityPublic); err != nil {
				return fmt.Errorf("field[%d]: %w", key, err)
			}
		}
	}

	for idx, ov := range file.Overwrites {
		if err := ov.Validate(); err != nil {
			return fmt.Errorf("overwrite[%d]: %w", idx, err)
		}
	}

	return nil
}

func (file Config) DynamicRolesEnabled() bool {
	if file.EnableDynamicRoles == nil {
		return false
	}

	return *file.EnableDynamicRoles
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
