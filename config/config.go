package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	exposingPort       = "EXPOSING_PORT"
	fileStorage        = "FILE_STORAGE"
	environment        = "ENVIRONMENT"
	featureFlag        = "FEATURE_FLAG"
	dbConnectionString = "DB_CONNECTION_STRING"
	redisAddr          = "REDIS_ADDR"
	redisPassword      = "REDIS_PASSWORD"
	redisDB            = "REDIS_DB"
)

const (
	// DevelopmentEnv ...
	DevelopmentEnv = "development"
	// ProductionEnv ...
	ProductionEnv = "production"
)

// Config contains application configuration
type Config struct {
	DBConnectionString string
	ExposingPort       string
	FileStorage        string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	FeatureFlag        string
	IsDevelopment      bool
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	redisDBi, err := strconv.Atoi(getEnvOrDefault(redisDB, "0"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis db: %v", err)
	}

	// default configuration
	config := &Config{
		DBConnectionString: getEnvOrDefault(dbConnectionString, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		ExposingPort:       getEnvOrDefault(exposingPort, "8089"),
		FileStorage:        getEnvOrDefault(fileStorage, "request.log"),
		RedisAddr:          getEnvOrDefault(redisAddr, "localhost:6379"),
		RedisPassword:      getEnvOrDefault(redisPassword, ""),
		RedisDB:            redisDBi,
		FeatureFlag:        getEnvOrDefault(featureFlag, "development"),
	}

	if env := os.Getenv(environment); env == DevelopmentEnv {
		os.Setenv(featureFlag, config.FeatureFlag)
		config.IsDevelopment = true
	}

	return config, nil
}
