package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	Application ApplicationConfig
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
}

type ApplicationConfig struct {
	Name        string
	Env         string
	FrontendURL string
}

type JWTConfig struct {
	SecretKey string
	Expiry    time.Duration
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLS          TLSServerConfig
}

type TLSServerConfig struct {
	Mode     string
	CertFile string
	KeyFile  string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      string

	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        time.Duration
}

func Load() *Configuration {
	config := &Configuration{
		Application: ApplicationConfig{
			Name:        GetEnv("APP_NAME", "social"),
			Env:         GetEnv("APP_ENV", "development"),
			FrontendURL: GetEnv("APP_FRONTEND", "http://localhost:5431"),
		},
		Server: ServerConfig{
			Port:         GetEnv("SERVER_PORT", "8080"),
			ReadTimeout:  GetEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: GetEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			TLS: TLSServerConfig{
				Mode:     GetEnv("SERVER_TLS_MODE", "off"),
				CertFile: GetEnv("SERVER_TLS_CERT_FILE", ""),
				KeyFile:  GetEnv("SERVER_TLS_KEY_FILE", ""),
			},
		},
		Database: DatabaseConfig{
			Host:               GetEnv("DB_HOST", "localhost"),
			Port:               GetEnv("DB_PORT", "5432"),
			User:               GetEnv("DB_USER", "admin"),
			Password:           GetEnv("DB_PASS", "password"),
			DatabaseName:       GetEnv("DB_NAME", "social"),
			SSLMode:            GetEnv("DB_SSL_MODE", "disable"),
			MaxOpenConnections: GetEnvInt("DB_MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConnections: GetEnvInt("DB_MAX_IDLE_CONNECTIONS", 10),
			MaxIdleTime:        GetEnvDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey: GetEnv("JWT_SECRET_KEY", "change-me-in-production"),
			Expiry:    GetEnvDuration("JWT_EXPIRY", 24*time.Hour),
		},
	}

	return config
}

func (c DatabaseConfig) DatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DatabaseName,
		c.SSLMode,
	)
}

func GetEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	val := GetEnv(key, "")
	if val == "" {
		return fallback
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return d
}
