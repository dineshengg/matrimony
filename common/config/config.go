package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	RedisHost  string `json:"redis_host"`
	RedisPort  int    `json:"redis_port"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig() (*Config, error) {
	// Define a flag for the config file path
	configFilePath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	// Open the JSON file
	file, err := os.Open(*configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// Decode the JSON file into the Config struct
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return &config, nil
}
