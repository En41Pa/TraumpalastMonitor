package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DiscordWebhookURL     string `json:"discord_webhook_url"`
	TraumpalastURL        string `json:"traumpalast_url"`
	CurrentNewestDate     string `json:"current_newest_date"`
	CheckIntervalMinutes  int    `json:"check_interval_minutes"`
	RequestTimeoutSeconds int    `json:"request_timeout_seconds"`
	UserAgent            string `json:"user_agent"`
}

func Load(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %v", filepath, err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Validate required fields
	if config.TraumpalastURL == "" {
		return nil, fmt.Errorf("traumpalast_url is required")
	}
	if config.CurrentNewestDate == "" {
		return nil, fmt.Errorf("current_newest_date is required")
	}
	if config.CheckIntervalMinutes <= 0 {
		config.CheckIntervalMinutes = 30 // default
	}
	if config.RequestTimeoutSeconds <= 0 {
		config.RequestTimeoutSeconds = 30 // default
	}

	return &config, nil
}
