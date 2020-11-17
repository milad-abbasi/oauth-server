package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Env file does not exist")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvWithDefault(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}

func MustGetEnv(key string) string {
	val, exists := os.LookupEnv(key)

	if !exists {
		panic(fmt.Sprintf("key \"%s\" does not exist in env", key))
	}

	return val
}
