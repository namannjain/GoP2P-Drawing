package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	OAuth    OAuthConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	OAuthRedirectURL   string
}

// Load reads the environment variables and returns a Config struct
func Load() (*Config, error) {
	// // Load .env file if it exists
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading the ENV file")
	// 	panic(err)
	// }

	config := &Config{}

	// Server Configuration
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Host = getEnv("SERVER_HOST", "localhost")

	// Database Configuration
	config.Database.Host = getEnv("DB_HOST", "postgres")
	config.Database.Port = getEnvAsInt("DB_PORT", 5432)
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "drawingapp")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// AWS Configuration
	config.AWS.Region = getEnv("AWS_REGION", "")
	config.AWS.AccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	config.AWS.SecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")
	config.AWS.BucketName = getEnv("AWS_BUCKET_NAME", "p2pdrawing")

	// OAuth Configuration
	config.OAuth.GoogleClientID = getEnv("GOOGLE_CLIENT_ID", "")
	config.OAuth.GoogleClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
	config.OAuth.OAuthRedirectURL = getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/auth/google/callback")

	fmt.Println(config)

	// Validate required configurations
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.OAuth.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if c.OAuth.GoogleClientSecret == "" {
		return fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if c.AWS.AccessKeyID == "" {
		return fmt.Errorf("AWS_ACCESS_KEY_ID is required")
	}
	if c.AWS.SecretAccessKey == "" {
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY is required")
	}
	return nil
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDatabaseURL returns the formatted database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
