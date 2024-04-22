// config/config.go

package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
    err := godotenv.Load()
    if err != nil {
        return fmt.Errorf("error loading .env file: %w", err)
    }
    return nil
}

// GetEnvVar returns the value of the environment variable with the given key
func GetEnvVar(key string) string {
    return os.Getenv(key)
}
