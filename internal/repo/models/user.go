package models

type Address struct {
	UserID   string `mapstructure:"user_id"`
	CityCode string `mapstructure:"city_code"`
	CityName string `mapstructure:"city_name"`
	Street   string `mapstructure:"street"`
	Extra    string `mapstructure:"extra"`
}

type EMail struct {
	UserID   string `mapstructure:"user_id"`
	Address  string `mapstructure:"address"`
	Verified bool   `mapstructure:"verified"`
}

type PhoneNumber struct {
	UserID      string `mapstructure:"user_id"`
	PhoneNumber string `mapstructure:"phone_number"`
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

type Group struct {
	ID          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

type GroupMembership struct {
	UserID  string `mapstructure:"user_id"`
	GroupID string `mapstructure:"group_id"`
}
