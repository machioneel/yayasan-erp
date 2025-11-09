package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Upload   UploadConfig
	Email    EmailConfig
	App      AppConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	AppName     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	Origins          []string
	AllowCredentials bool
}

type UploadConfig struct {
	MaxSize        int64
	Path           string
	AllowedTypes   []string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

type AppConfig struct {
	EnableAuditLog    bool
	EnableMultiBranch bool
	EnableMultiCurrency bool
	DefaultCurrency   string
	DefaultLanguage   string
	DefaultTimezone   string
	FiscalYearStart   int
	RateLimitEnabled  bool
	RateLimitPerMin   int
	DefaultPageSize   int
	MaxPageSize       int
}

var AppConfig *Config

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	jwtRefreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		jwtRefreshExpiry = 168 * time.Hour
	}

	config := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENV", "development"),
			AppName:     getEnv("APP_NAME", "Yayasan ERP"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "yayasan_erp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "change-this-secret-key"),
			Expiry:        jwtExpiry,
			RefreshExpiry: jwtRefreshExpiry,
		},
		CORS: CORSConfig{
			Origins:          strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
		},
		Upload: UploadConfig{
			MaxSize:      getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB default
			Path:         getEnv("UPLOAD_PATH", "./uploads"),
			AllowedTypes: strings.Split(getEnv("ALLOWED_FILE_TYPES", "image/jpeg,image/png,application/pdf"), ","),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPFrom:     getEnv("SMTP_FROM", "noreply@yayasan.org"),
		},
		App: AppConfig{
			EnableAuditLog:      getEnvAsBool("ENABLE_AUDIT_LOG", true),
			EnableMultiBranch:   getEnvAsBool("ENABLE_MULTI_BRANCH", true),
			EnableMultiCurrency: getEnvAsBool("ENABLE_MULTI_CURRENCY", false),
			DefaultCurrency:     getEnv("DEFAULT_CURRENCY", "IDR"),
			DefaultLanguage:     getEnv("DEFAULT_LANGUAGE", "id"),
			DefaultTimezone:     getEnv("DEFAULT_TIMEZONE", "Asia/Jakarta"),
			FiscalYearStart:     getEnvAsInt("FISCAL_YEAR_START_MONTH", 1),
			RateLimitEnabled:    getEnvAsBool("RATE_LIMIT_ENABLED", true),
			RateLimitPerMin:     getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			DefaultPageSize:     getEnvAsInt("DEFAULT_PAGE_SIZE", 20),
			MaxPageSize:         getEnvAsInt("MAX_PAGE_SIZE", 100),
		},
	}

	// Validate required fields
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	if config.JWT.Secret == "change-this-secret-key" && config.Server.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be changed in production")
	}

	AppConfig = config
	return config, nil
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
		c.Timezone,
	)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// IsDevelopment checks if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction checks if environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
