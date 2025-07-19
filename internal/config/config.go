package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	File     FileConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	ReadTimeout  int
	WriteTimeout int
	EnableLogs   bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int
}

type JWTConfig struct {
	Secret    string
	ExpiresIn int // hours
}

type FileConfig struct {
	UploadPath string
	MaxSize    int64 // bytes
	AllowedTypes []string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 10),
			EnableLogs:   getEnvAsBool("ENABLE_LOGS", true),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "messenger"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			Database: getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		File: FileConfig{
			UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
			MaxSize:    getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
			AllowedTypes: []string{
				"image/jpeg", "image/png", "image/gif", "image/webp",
				"video/mp4", "video/avi", "video/mov", "video/webm",
				"audio/mp3", "audio/wav", "audio/ogg", "audio/m4a",
				"application/pdf", "application/doc", "application/docx",
				"application/zip", "application/rar", "text/plain",
			},
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}