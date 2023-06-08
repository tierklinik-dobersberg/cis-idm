package models

type Address struct {
	UserID   string `mapstructure:"user_id"`
	CityCode string `mapstructure:"city_code"`
	CityName string `mapstructure:"city_name"`
	Street   string `mapstructure:"street"`
	Extra    string `mapstructure:"extra"`
}

type EMail struct {
	UserID   string   `mapstructure:"user_id"`
	Address  string   `mapstructure:"address"`
	Verified bool     `mapstructure:"verified"`
	Primary  bool     `mapstructure:"primary"`
	Tags     []string `mapstructure:"tags"`
}

type PhoneNumber struct {
	UserID      string   `mapstructure:"user_id"`
	PhoneNumber string   `mapstructure:"phone_number"`
	Primary     bool     `mapstructure:"primary"`
	Tags        []string `mapstructure:"tags"`
}

type User struct {
	ID          string `mapstructure:"id"`
	Username    string `mapstructure:"username"`
	DisplayName string `mapstructure:"display_name"`
	FirstName   string `mapstructure:"first_name"`
	LastName    string `mapstructure:"last_name"`
	Extra       []byte `mapstructure:"extra"`
	Avatar      string `mapstructure:"avatar"`
	Password    string `mapstructure:"password"`
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
	TokenID  string `mapstructure:"token_id"`
	UserID   string `mapstructure:"user_id"`
	IssuedAt int64  `mapstructure:"issued_at"`
	ExiresAt int64  `mapstructure:"expires_at"`
}
