package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	DatabaseURL   string
	JWTSecret     string
	HTTPPort      string
	ReferralBonus int64
	RefereeBonus  int64
}

// Load reads configuration from environment variables
//
//nolint:unused // Will be used when main.go is implemented
func load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/user_management?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev_secret_change_in_production"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
	}
	var err error

	// Parse referral bonus from string to int64
	ReferralBonusStr := getEnv("REFERRAL_BONUS", "100")
	cfg.ReferralBonus, err = strconv.ParseInt(ReferralBonusStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid REFERRAL_BONUS: %v", err)
	}

	// Parse referee bonus from string to int64
	RefereeBonusStr := getEnv("REFEREE_BONUS", "50")
	cfg.RefereeBonus, err = strconv.ParseInt(RefereeBonusStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid REFEREE_BONUS %v", err)
	}
	return cfg, nil
}

// getEnv returns environment variable or default value
//
// nolint:unused // Will be used when main.go is implemented
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
