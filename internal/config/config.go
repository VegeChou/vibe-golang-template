package config

import "os"

type Config struct {
	HTTPAddr string
	I18NFile string
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		I18NFile: getEnv("I18N_FILE", "configs/i18n.json"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
