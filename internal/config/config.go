package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	MCP      MCPConfig
	API      APIConfig
	Security SecurityConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// MCPConfig holds MCP server configuration
type MCPConfig struct {
	HTTPPort int
	LogLevel string
}

// APIConfig holds REST API configuration
type APIConfig struct {
	Port   int
	WSPort int
}

// SecurityConfig holds authentication and authorization configuration
type SecurityConfig struct {
	OIDCIssuerURL string
	OIDCClientID  string
	OPAURL        string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://synthesis:synthesis_dev_password@localhost:5432/synthesis?sslmode=disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 300),
		},
		MCP: MCPConfig{
			HTTPPort: getEnvAsInt("MCP_HTTP_PORT", 8081),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		API: APIConfig{
			Port:   getEnvAsInt("API_PORT", 8080),
			WSPort: getEnvAsInt("WS_PORT", 8082),
		},
		Security: SecurityConfig{
			OIDCIssuerURL: getEnv("OIDC_ISSUER_URL", ""),
			OIDCClientID:  getEnv("OIDC_CLIENT_ID", ""),
			OPAURL:        getEnv("OPA_URL", "http://localhost:8181"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate ensures required configuration is present
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
