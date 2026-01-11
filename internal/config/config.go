package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	ServerHost string
	ServerURL  string // Backend server URL for callbacks (e.g., http://api.domain.com or http://192.168.1.100:5000)
	ClientURL  string // Frontend client URL (for CORS)

	// Database
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string
	DatabaseURL      string

	// JWT
	JWTSecret string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// RabbitMQ
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string

	// Email
	EmailFrom    string
	EmailName    string // Custom sender name (e.g., "Zacode")
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string

	// Rate Limiting
	RateLimitEnabled bool
	RateLimitRPS     int // Requests per second
	RateLimitBurst   int // Burst size

	// Midtrans Payment Gateway
	MidtransServerKey string
	MidtransClientKey string

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	serverPort := getEnv("PORT", "5000")
	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverURL := getEnv("SERVER_URL", "") // Backend URL for callbacks
	// If SERVER_URL not set, construct from SERVER_HOST and PORT
	if serverURL == "" {
		if serverHost == "0.0.0.0" || serverHost == "" {
			serverURL = fmt.Sprintf("http://localhost:%s", serverPort)
		} else {
			serverURL = fmt.Sprintf("http://%s:%s", serverHost, serverPort)
		}
	}

	cfg := &Config{
		// Server
		ServerPort: serverPort,
		ServerHost: serverHost,
		ServerURL:  serverURL,
		ClientURL:  getEnv("CLIENT_URL", "http://localhost:3000"),

		// Database
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "yourapp"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),

		// Google OAuth
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// RabbitMQ
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),

		// Email
		EmailFrom:    getEnv("EMAIL_FROM", ""),
		EmailName:    getEnv("EMAIL_NAME", "Zacode"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),

		// Rate Limiting (default: enabled, 100 req/sec, burst 200)
		RateLimitEnabled: getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitRPS:     getEnvInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:   getEnvInt("RATE_LIMIT_BURST", 200),

		// Midtrans Payment Gateway
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-4zIt7djwCeRdMpgF4gXDjciC"),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),

		// Cloudinary
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", "dgmlqboeq"),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", "736499913818945"),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", "pfFz2h0qhf8qTIEGWEjQQbqsYWk"),
	}

	// Build database URL if not provided
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.PostgresHost,
			cfg.PostgresPort,
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresDB,
			cfg.PostgresSSLMode,
		)
	}

	// Validate required fields
	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-secret-key-change-in-production" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" || value == "yes" {
			return true
		}
		return false
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
