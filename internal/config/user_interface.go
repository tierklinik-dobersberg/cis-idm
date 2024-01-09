package config

import "fmt"

type UserInterface struct {
	// SiteName can be used to specify the name of the cisidm instance and will be displayed
	// at the login screen and throughout the user interface. This defaults to Example
	// so will likely want to set this field as well.
	SiteName string `json:"site_name" hcl:"site_name"`

	// SiteNameURL can be set to a URL that will be used to create a HTML link on the login
	// page.
	SiteNameURL string `json:"site_name_url" hcl:"site_name_url,optional"`

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

	// LogoURL may be set to a path or HTTP resource that should be displayed as the
	// application logo on the login screen.
	LogoURL string `json:"logo_url" hcl:"logo_url,optional"`

	// PublicURL defines the public URL at which cisidm is reachable from the outside.
	// This value MUST be set.
	PublicURL string `json:"public_url" hcl:"public_url"`
}

func (ui *UserInterface) ApplyDefaultsAndValidate() error {
	if ui == nil {
		return fmt.Errorf("missing configuration block")
	}

	if ui.PublicURL == "" {
		return fmt.Errorf("public_url is missing")
	}

	if ui.SiteNameURL == "" {
		ui.SiteNameURL = ui.PublicURL
	}

	if ui.LoginRedirectURL == "" {
		ui.LoginRedirectURL = fmt.Sprintf("%s/login?redirect=%%s", ui.PublicURL)
	}

	if ui.RefreshRedirectURL == "" {
		ui.RefreshRedirectURL = fmt.Sprintf("%s/refresh?redirect=%%s", ui.PublicURL)
	}

	if ui.PasswordResetURL == "" {
		ui.PasswordResetURL = fmt.Sprintf("%s/password/reset?token=%%s", ui.PublicURL)
	}

	if ui.VerifyMailURL == "" {
		ui.VerifyMailURL = fmt.Sprintf("%s/profile/verify-mail?token=%%s", ui.PublicURL)
	}

	if ui.RegistrationURL == "" {
		ui.RegistrationURL = fmt.Sprintf("%s/registration?token=%%s&mail=%%s&name=%%s", ui.PublicURL)
	}

	return nil
}
