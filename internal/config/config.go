package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// App-wide structured config
type AppConfig struct {
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
		Port string
	} `yaml:"app"`

	DatabaseURL string
	JWTSecret   string

	Redis struct {
		Address string `yaml:"address"`
		DB      int    `yaml:"db"`
	} `yaml:"redis"`
}

// Env + DB + JWT secrets config
type DBConfig struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	DatabaseURL      string
	JWTSecret        string
	JWTRefreshSecret string
	RedisURL         string
	WaafiMerchantUID string
	Env              string
}

// LoadEnv loads from .env file into OS env vars
func LoadEnv() {
	err := godotenv.Load("/app/.env") // absolute path for Docker
	if err != nil {
		log.Println("No .env file found; falling back to system env vars")
	}
}

// getEnv gets env var or fallback
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// LoadAppConfig loads YAML + overrides from env
func LoadAppConfig() (*AppConfig, error) {
	file, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	// override with env vars
	cfg.App.Port = getEnv("PORT", cfg.App.Port)
	cfg.DatabaseURL = getEnv("DATABASE_URL", cfg.DatabaseURL)
	cfg.JWTSecret = getEnv("JWT_SECRET", cfg.JWTSecret)

	return &cfg, nil
}

// LoadDBConfig returns DBConfig using only env vars
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Port:             getEnv("PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "kaabe_user"),
		DBPassword:       getEnv("DB_PASSWORD", "kaabe_password"),
		DBName:           getEnv("DB_NAME", "kaabe"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:        getEnv("JWT_SECRET", "default_jwt_secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "default_jwt_refresh_secret"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		WaafiMerchantUID: getEnv("WAAFI_MERCHANT_UID", ""),
		Env:              getEnv("ENV", "development"),
	}
}

// InitDB connects to PostgreSQL using database/sql
func InitDB(cfg *DBConfig) *sql.DB {
	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to the database successfully.")
	return db
}
