package config

import "os"

type Config struct {
	Port       string
	Secret     string
	Branch     string
	ApplyScript string
}

func Load() *Config {
	return &Config{
		Port:        getEnvOrDefault("PORT", "8050"),
		Secret:      os.Getenv("WEBHOOK_SECRET"),
		Branch:      getEnvOrDefault("WEBHOOK_BRANCH", "master"),
		ApplyScript: getEnvOrDefault("APPLY_SCRIPT", "/opt/sysopoly-infrastructure/apply.sh"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
