package model

// Config holds the application configuration
type Config struct {
	Port       int    `yaml:"port" flag:"p" env:"WEBSITE_PORT" default:"8080"`
	DBCS       string `yaml:"dbcs" flag:"dbcs" env:"WEBSITE_DATABASE_CONNECTION_STRING"`
	SigningKey string `yaml:"signingkey" flag:"s" env:"WEBSITE_SIGNING_KEY"`
	OID        struct {
		URL          string `yaml:"url" flag:"oidp" env:"WEBSITE_OID_PROVIDER"`
		ClientID     string `yaml:"client_id" flag:"cid" env:"WEBSITE_OID_CLIENT_ID"`
		ClientSecret string `yaml:"client_secret" flag:"csecret" env:"WEBSITE_OID_CLIENT_SECRET"`
		StateString  string `yaml:"statestring" flag:"state" env:"WEBSITE_OID_STATE_STRING"`
	} `yaml:"open_id"`
}

type UserInfo struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Verified bool     `json:"verified"`
	Roles    []string `json:"roles"`
}
