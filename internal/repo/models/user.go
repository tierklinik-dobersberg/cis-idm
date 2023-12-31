package models

type Address struct {
	ID       string `mapstructure:"id"`
	UserID   string `mapstructure:"user_id"`
	CityCode string `mapstructure:"city_code"`
	CityName string `mapstructure:"city_name"`
	Street   string `mapstructure:"street"`
	Extra    string `mapstructure:"extra"`
}

type EMail struct {
	ID       string   `mapstructure:"id"`
	UserID   string   `mapstructure:"user_id"`
	Address  string   `mapstructure:"address"`
	Verified bool     `mapstructure:"verified"`
	Primary  bool     `mapstructure:"is_primary"`
	Tags     []string `mapstructure:"tags"`
}

type PhoneNumber struct {
	ID          string   `mapstructure:"id"`
	UserID      string   `mapstructure:"user_id"`
	PhoneNumber string   `mapstructure:"phone_number"`
	Primary     bool     `mapstructure:"is_primary"`
	Verified    bool     `mapstructure:"verified"`
	Tags        []string `mapstructure:"tags"`
}

type User struct {
	ID          string `mapstructure:"id"`
	Username    string `mapstructure:"username"`
	DisplayName string `mapstructure:"display_name"`
	FirstName   string `mapstructure:"first_name"`
	LastName    string `mapstructure:"last_name"`
	Extra       string `mapstructure:"extra"`
	Avatar      string `mapstructure:"avatar"`
	Password    string `mapstructure:"password"`
	Birthday    string `mapstructure:"birthday"`
	TOTPSecret  string `mapstructure:"totp_secret"`
}

type Role struct {
	ID              string `mapstructure:"id"`
	Name            string `mapstructure:"name"`
	Description     string `mapstructure:"description"`
	DeleteProtected bool   `mapstructure:"delete_protected"`
}

type RoleAssignment struct {
	UserID string `mapstructure:"user_id"`
	RoleID string `mapstructure:"role_id"`
}

type RejectedToken struct {
	TokenID   string `mapstructure:"token_id"`
	UserID    string `mapstructure:"user_id"`
	IssuedAt  int64  `mapstructure:"issued_at"`
	ExpiresAt int64  `mapstructure:"expires_at"`
}

type Feature struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

type UserEnabledFeature struct {
	FeatureName string `mapstructure:"feature_name"`
	UserID      string `mapstructure:"user_id"`
}

type RegistrationToken struct {
	InitialRoles string `mapstructure:"initial_roles"`
	Token        string `mapstructure:"token"`
	Expires      *int64 `mapstructure:"expires"`
	AllowedUsage *int64 `mapstructure:"allowed_usage"`
	CreatedBy    string `mapstructure:"created_by"`
	CreatedAt    int64  `mapstructure:"created_at"`
}

type Passkey struct {
	ID           string `mapstructure:"id"`
	UserID       string `mapstructure:"user_id"`
	Cred         string `mapstructure:"cred"`
	CredType     string `mapstructure:"cred_type"`
	ClientName   string `mapstructure:"client_name"`
	ClientOS     string `mapstructure:"client_os"`
	ClientDevice string `mapstructure:"client_device"`
}
