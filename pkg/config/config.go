package config

import (
	"fmt"
	"os"
)

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

func MustGet(key string) string {
	val, exists := os.LookupEnv(key)

	if !exists {
		panic(fmt.Sprintf("key \"%s\" does not exist in env", key))
	}

	return val
}
