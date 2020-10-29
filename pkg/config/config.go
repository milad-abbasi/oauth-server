package config

import "os"

func Get(key string) string {
	return os.Getenv(key)
}

func GetWithDefault(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
