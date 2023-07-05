package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/go-multierror"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slices"
)

type File struct {
	ForwardAuth []*ForwardAuthEntry `json:"forwardAuth" yaml:"forwardAuth"`

	Audience               string       `json:"audience"`
	JWTSecret              string       `json:"jwtSecret"`
	DatabaseURL            string       `json:"rqliteURL"`
	SecureCookie           bool         `json:"secureCookie"`
	AccessTokenTTL         JSONDuration `json:"accessTokenTTL"`
	RefreshTokenTTL        JSONDuration `json:"refreshTokenTTL"`
	AccessTokenCookieName  string       `json:"accessTokenCookieName"`
	RefreshTokenCookieName string       `json:"refreshTokenCookieName"`
	BootstrapRoles         []string     `json:"bootstrapRoles"`
	AllowedDomainRedirects []string     `json:"allowedRedirects"`
	FeatureSet             []Feature    `json:"features"`
	PublicListenAddr       string       `json:"publicListener"`
	AdminListenAddr        string       `json:"adminListener"`
	AllowedOrigins         []string     `json:"allowedOrigins"`
	PublicURL              string       `json:"publicURL"`
	StaticFiles            string       `json:"staticFiles"`

	// Exposed via /config endpoint
	RegistrationRequiresToken bool   `json:"registrationRequiresToken"`
	Domain                    string `json:"domain"`
	LoginRedirectURL          string `json:"loginURL"`
	RefreshRedirectURL        string `json:"refreshURL"`
	SiteName                  string `json:"siteName"`
	SiteNameURL               string `json:"siteNameUrl"`

	// built from FeatureSet
	FeatureMap map[Feature]bool `json:"-"`
}

func LoadFile(path string) (*File, error) {
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

	var f File
	if err := json.Unmarshal(content, &f); err != nil {
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

func (file *File) applyDefaults() error {
	if file.PublicURL == "" {
		return fmt.Errorf("publicURL must be set")
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
		file.AdminListenAddr = "localhost:8081"
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

	if file.Audience == "" || file.Domain == "" {
		return fmt.Errorf("missing domain and audience")
	}

	return nil
}

func (file File) AuthRequiredForURL(method string, url string) (bool, error) {
	merr := new(multierror.Error)

	for _, fae := range file.ForwardAuth {
		matches, err := fae.Matches(method, url)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("invalid regex %q: %w", fae.URL, err))
			continue
		}

		if matches {
			if fae.IsRequired() {
				return true, merr.ErrorOrNil()
			}
		}
	}

	return false, merr.ErrorOrNil()
}

type Config struct {
	File

	ConfigFilePath string `env:"CONFIG_FILE,default=/etc/cisidm/config.yml"`
}

// FromEnvironment returns a Config object parsed from environment variables.
func FromEnvironment(ctx context.Context) (cfg Config, err error) {
	l := envconfig.PrefixLookuper("IDM_", envconfig.OsLookuper())

	if err := envconfig.ProcessWith(ctx, &cfg, l); err != nil {
		return cfg, err
	}

	cfg.parseFeatureSet()

	if cfg.ConfigFilePath != "" {
		parsedFile, err := LoadFile(cfg.ConfigFilePath)
		if err != nil {
			return cfg, fmt.Errorf("failed to parse config file at %q: %w", cfg.ConfigFilePath, err)
		}

		cfg.File = *parsedFile
	}

	return cfg, nil
}

func (file *File) parseFeatureSet() error {
	defaultValue := slices.Contains(file.FeatureSet, FeatureAll)

	file.FeatureMap = make(map[Feature]bool)
	for _, feat := range AllFeatures {
		file.FeatureMap[feat] = defaultValue
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

		file.FeatureMap[feat] = allowed
	}

	return nil
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
