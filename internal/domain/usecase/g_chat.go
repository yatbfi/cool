package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GChat defines the Google Chat usecase interface
type GChat interface {
	// SendMessage sends a message to Google Chat webhook
	SendMessage(ctx context.Context, webhookURL string, message string) error
}

// gchatUsecase implements GChat interface
type gchatUsecase struct {
	httpClient *http.Client
}

// NewGChatUsecase creates a new Google Chat usecase
func NewGChatUsecase() GChat {
	return &gchatUsecase{
		httpClient: &http.Client{},
	}
}

// SendMessage sends a message to Google Chat webhook
func (u *gchatUsecase) SendMessage(ctx context.Context, webhookURL string, message string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	payload := map[string]interface{}{
		"text": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
