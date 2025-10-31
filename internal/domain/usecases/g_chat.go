package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/pkg/logger"
)

type GChat interface {
	SendToChannel(ctx context.Context, name string, message string) error
}

type GChatUsecase struct {
	cfg *config.Config
}

var _ GChat = (*GChatUsecase)(nil)

func NewGChatUsecase(cfg *config.Config) *GChatUsecase {
	return &GChatUsecase{cfg: cfg}
}

func (u *GChatUsecase) SendToChannel(ctx context.Context, name string, message string) error {
	webhookURL, err := u.getWebhookUrl(name)
	if err != nil {
		return fmt.Errorf("SendToChannel: %w", err)
	}

	// --- Enrich message with user info ---
	if u.cfg.UserName != "" || u.cfg.UserEmail != "" {
		userInfo := fmt.Sprintf("\n\n👤 *%s* <%s>", u.cfg.UserName, u.cfg.UserEmail)
		message = fmt.Sprintf("%s%s", message, userInfo)
	}

	payload := map[string]string{
		"text": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.Errorf("failed to close response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-success response from Google Chat: %s", resp.Status)
	}

	return nil
}

func (u *GChatUsecase) getWebhookUrl(name string) (string, error) {
	switch name {
	case "lora_code_review":
		return u.cfg.GChatReviewWebhookURL, nil
	default:
		return "", fmt.Errorf("no webhook URL configured for channel: %s", name)
	}
}
