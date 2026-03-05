package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration values loaded from environment
// variables. Using a struct makes it easy to pass config around and test
// components by injecting a fake config.

type Config struct {
    Port           string
    DBUrl          string
    JWTSecret      string
    KafkaBroker    string
    RedisAddr      string
    RedisPassword  string
    SMTPHost       string
    SMTPPort       int
    SMTPUser       string
    SMTPPassword   string
    OpenAIKey      string
}

// Load reads environment variables (with support for a .env file) and
// returns a complete Config. It will log.Fatal if required variables are
// missing or malformed.

func Load() *Config {
    // allow loading from .env for local development
    _ = godotenv.Load()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        smtpPort = 587
    }

    cfg := &Config{
        Port:          port,
        DBUrl:         os.Getenv("DATABASE_URL"),
        JWTSecret:     os.Getenv("JWT_SECRET"),
        KafkaBroker:   os.Getenv("KAFKA_BROKER"),
        RedisAddr:     os.Getenv("REDIS_ADDR"),
        RedisPassword: os.Getenv("REDIS_PASSWORD"),
        SMTPHost:      os.Getenv("SMTP_HOST"),
        SMTPPort:      smtpPort,
        SMTPUser:      os.Getenv("SMTP_USER"),
        SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
        OpenAIKey:     os.Getenv("OPENAI_KEY"),
    }

    if cfg.DBUrl == "" || cfg.JWTSecret == "" || cfg.KafkaBroker == "" {
        log.Fatal("required environment variables are not set")
    }

    return cfg
}
