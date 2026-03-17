package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config struct definition
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Environment string
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host          string
	Port          int
	Name          string
	User          string
	Password      string
	SSLMode       string
	RunMigrations bool
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

func LoadConfig() *Config {
	// Set defaults
	viper.SetDefault("Server.Host", "localhost")
	viper.SetDefault("Server.Port", 8080)
	viper.SetDefault("Database.Host", "localhost")
	viper.SetDefault("Database.Port", 5432)
	viper.SetDefault("Database.Name", "payables_db")
	viper.SetDefault("Database.User", "payables_user")
	viper.SetDefault("Database.Password", "payables_password")
	viper.SetDefault("Database.SSLMode", "disable")
	viper.SetDefault("JWT.Secret", "your_secret_key_here")
	viper.SetDefault("JWT.Expiration", "24h")
	viper.SetDefault("Environment", "development")

	// Enable environment variables
	viper.AutomaticEnv()

	// Replace dots with underscores for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	return &config
}
