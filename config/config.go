package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

// Config holds environment and local user configuration.
type Config struct {
	GChatReviewWebhookURL string `envconfig:"G_CHAT_REVIEW_WEBHOOK_URL" default:"https://chat.googleapis.com/v1/spaces/AAQAzjrjbvI/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Jk5CAkN3LaeFSAtjeRUmeJHY_np-HQ89eI-5cmYEQrI"`
	GChatCollabWebhookURL string `envconfig:"G_CHAT_COLLAB_WEBHOOK_URL" default:"https://chat.googleapis.com/v1/spaces/AAQAzjrjbvI/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Jk5CAkN3LaeFSAtjeRUmeJHY_np-HQ89eI-5cmYEQrI"`

	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

var cached *Config

// GetConfig loads configuration from environment variables
// and merges user info from ~/.cool-cli/config.json if present.
func GetConfig() *Config {
	if cached != nil {
		return cached
	}

	cfg := &Config{}
	envconfig.MustProcess("", cfg)

	if local, err := loadLocalConfig(); err == nil {
		if local.UserName != "" {
			cfg.UserName = local.UserName
		}
		if local.UserEmail != "" {
			cfg.UserEmail = local.UserEmail
		}
	}

	cached = cfg
	return cached
}

// SaveLocalConfig stores UserName and UserEmail in ~/.cool-cli/config.json.
func SaveLocalConfig(c *Config) error {
	path := getLocalConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(struct {
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
	}{
		UserName:  c.UserName,
		UserEmail: c.UserEmail,
	}, "", "  ")
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
