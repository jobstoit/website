package model

// Config holds the application configuration
type Config struct {
	Port int    `yaml:"port" flag:"p" env:"WEBSITE_PORT" default:"8080"`
	DBCS string `yaml:"dbcs" flag:"dbcs" env:"WEBSITE_DATABASE_CONNECTION_STRING"`
	OIDP string `yaml:"oidp" flag:"oidp" env:"WEBSITE_OPEN_ID_PROVIDER"`
}

type UserInfo struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Verified bool     `json:"verified"`
	Roles    []string `json:"roles"`
}
