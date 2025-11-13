package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds environment and local user configuration.
type Config struct {
	GChatReviewWebhookURL string `json:"gchat_review_webhook_url,omitempty"`
	GChatCollabWebhookURL string `json:"gchat_collab_webhook_url,omitempty"`

	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

var cached *Config

// GetConfig loads configuration from ~/.cool-cli/config.json.
func GetConfig() *Config {
	if cached != nil {
		return cached
	}

	cfg := &Config{}

	if local, err := loadLocalConfig(); err == nil {
		cfg.UserName = local.UserName
		cfg.UserEmail = local.UserEmail
		cfg.GChatReviewWebhookURL = local.GChatReviewWebhookURL
		cfg.GChatCollabWebhookURL = local.GChatCollabWebhookURL
	}

	cached = cfg
	return cached
}

// SaveLocalConfig stores configuration in ~/.cool-cli/config.json.
func SaveLocalConfig(c *Config) error {
	path := getLocalConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	if err = os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func loadLocalConfig() (*Config, error) {
	path := getLocalConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("read local config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse local config: %w", err)
	}
	return &cfg, nil
}

func getLocalConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(home, ".cool-cli", "config.json")
}
