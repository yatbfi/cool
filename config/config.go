package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

// Config defines configuration loaded from environment and local setup.
type Config struct {
	GChatReviewWebhookURL string `envconfig:"G_CHAT_REVIEW_WEBHOOK_URL" default:"https://chat.googleapis.com/v1/spaces/AAQAzjrjbvI/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Jk5CAkN3LaeFSAtjeRUmeJHY_np-HQ89eI-5cmYEQrI"`

	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

var cfg *Config

// GetConfig loads config from environment and merges with ~/.cool-cli/config.json (if present).
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}

	// Step 1: load from env
	env := &Config{}
	envconfig.MustProcess("", env)

	// Step 2: try load from ~/.cool-cli/config.json
	fileCfg, _ := loadLocalConfig()

	// Step 3: merge (file overrides env for UserName/UserEmail)
	if fileCfg.UserName != "" {
		env.UserName = fileCfg.UserName
	}
	if fileCfg.UserEmail != "" {
		env.UserEmail = fileCfg.UserEmail
	}

	cfg = env
	return cfg
}

// SaveLocalConfig saves the current UserName/UserEmail into ~/.cool-cli/config.json
func SaveLocalConfig(c *Config) error {
	path := getLocalConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
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

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse local config: %w", err)
	}
	return &c, nil
}

func getLocalConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cool-cli", "config.json")
}
