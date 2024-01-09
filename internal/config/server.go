package config

import "fmt"

type Server struct {
	// SecureCookie defines whether or not cookies should be set with the
	// Secure attribute. If left empty, SecureCookie will be automatically
	// set depending on the PublicURL field.
	SecureCookie *bool `json:"secure_cookies" hcl:"secure_cookies,optional"`

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

	// AllowedOrigins configures a list of allowed origins for Cross-Origin-Requests.
	// This defaults to the PublicURL as well as http(s)://{{ Domain }}
	AllowedOrigins []string `json:"allowed_origins" hcl:"allowed_origins,optional"`

	// Domain is the parent domain for which cisidm handles authentication. If you
	// have multiple sub-domains hosting your services you want to set this to the
	// parent domain.
	//
	// I.e. if cisidm is running on account.example.com and you have services on
	// foo.example.com and bar.example.com you want to set the Domain field to "example.com"
	Domain string `json:"domain" hcl:"domain"`

	// TrustedNetworks is a list of CIDR network addresses that are considered
	// trusted. Any X-Forwareded-For header from these networks will be parsed
	// and applied.
	TrustedNetworks []string `json:"trusted_networks" hcl:"trusted_networks,optional"`

	// AllowedDomainRedirects is a list of domain names to which cisidm will allow
	// redirection after login/refresh.
	AllowedDomainRedirects []string `json:"allowed_redirects" hcl:"allowed_redirects,optional"`
}

func (file *Server) ApplyDefaultsAndValidate(secure bool) error {
	if file == nil {
		return fmt.Errorf("missing configuration block")
	}

	if file.PublicListenAddr == "" {
		file.PublicListenAddr = ":8080"
	}

	if file.AdminListenAddr == "" {
		file.AdminListenAddr = ":8081"
	}

	if file.SecureCookie == nil {
		file.SecureCookie = &secure
	}

	return nil
}
