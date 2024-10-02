package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// config holds all the configuration for the application
type config struct {
   ServerAddress string
   DatabasePath  string
   JWTSecret     string
}

// reads the configuration from the environment variables
func Load() (*config, error) {
   // load .env file if it exists
   _ = godotenv.Load()

   cfg := &config{
      ServerAddress: getEnv("SERVER_ADDRESS", ":7438"),
      DatabasePath: getEnv("DATABASE_PATH", "./taskmaster.db"),
      JWTSecret: getEnv("JWT_SECRET", ""),
   }

   if cfg.JWTSecret == "" {
      return nil, fmt.Errorf("JWT_SECRET must be set")
   }

   return cfg, nil
}

// retrieves an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
   value := os.Getenv(key)
   if value == "" {
      return defaultValue
   }
   return value
}
