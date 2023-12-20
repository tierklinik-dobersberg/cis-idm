package models

type Webpush struct {
	ID        string `mapstructure:"id"`
	UserID    string `mapstructure:"user_id"`
	UserAgent string `mapstructure:"user_agent"`
	Endpoint  string `mapstructure:"endpoint"`
	Auth      string `mapstructure:"auth"`
	Key       string `mapstructure:"key"`
	TokenID   string `mapstructure:"token_id"`
}
