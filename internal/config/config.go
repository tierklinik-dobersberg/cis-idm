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
	Directories      []string `json:"directories" hcl:"directories,optional"`
	ForwardAuthQuery string   `json:"forward_auth_query" hcl:"forward_auth_query,optional"`
	Debug            bool     `json:"debug" hcl:"debug,optional"`

	Policies []Policy `json:"policy" hcl:"policy,block"`
}

func (cfg *PolicyConfig) ApplyDefaultsAndValidate() error {
	if cfg == nil {
		return nil
	}

	if cfg.ForwardAuthQuery == "" {
		cfg.ForwardAuthQuery = "data.cisidm.forward_auth"
	}

	return nil
}

type Config struct {
	// LogLevel defines the log level to use.
	LogLevel string `json:"log_level" hcl:"log_level,optional"`

	// PolicyConfig holds the configuration for rego policies.
	PolicyConfig PolicyConfig `json:"policy" hcl:"policies,block"`

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
	Overwrites []Overwrite `json:"overwrites" hcl:"overwrites,block"`

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

	// RegistrationRequiresToken defines whether or not users are allowed to sign
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
				return fmt.Errorf("extraData: #%d: %w", key, err)
			}
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
