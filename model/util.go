package model

// Config holds the application configuration
type Config struct {
	Port int    `yaml:"port" flag:"p" env:"WEBSITE_PORT" default:"8080"`
	DBCS string `yaml:"dbcs" flag:"dbcs" env:"WEBSITE_DATABASE_CONNECTION_STRING"`
}
