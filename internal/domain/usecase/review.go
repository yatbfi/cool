package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/domain/entity"
	"github.com/yatbfi/cool/internal/domain/repository"
)

// ReviewRequest represents a review request input
type ReviewRequest struct {
	Title       string
	Description string
	Priority    string
	ReviewLinks []string
	JiraLinks   []string
}

// ReviewHistoryEntry represents a review history entry (alias from entity)
type ReviewHistoryEntry = entity.ReviewHistoryEntry

// HistoryFilter defines filter type for history queries
type HistoryFilter int

const (
	HistoryFilterAll HistoryFilter = iota
	HistoryFilterPending
	HistoryFilterCompleted
)

// Review defines the review usecase interface
type Review interface {
	// SubmitReviewRequest submits a new review request to tech lead
	SubmitReviewRequest(ctx context.Context, req *ReviewRequest, withSend bool) (*ReviewHistoryEntry, error)

	// FormatReviewRequestMessage formats review request for preview/sending
	FormatReviewRequestMessage(entry *ReviewHistoryEntry) string

	// GetHistories retrieves review histories with optional filter
	GetHistories(ctx context.Context, filter HistoryFilter) ([]*ReviewHistoryEntry, error)

	// GetHistoryByID retrieves a specific history by ID
	GetHistoryByID(ctx context.Context, id string) (*ReviewHistoryEntry, error)

	// SubmitToCollaboration forwards a review request to collaboration channel (head architect)
	SubmitToCollaboration(ctx context.Context, historyID string) error

	// SendToGChat sends a message to Google Chat webhook
	SendToGChat(ctx context.Context, webhookURL string, message string) error
}

// reviewUsecase implements Review interface
type reviewUsecase struct {
	historyRepo repository.ReviewHistoryRepository
	gchatUc     GChat
}

// NewReviewUsecase creates a new review usecase
func NewReviewUsecase(historyRepo repository.ReviewHistoryRepository, gchatUc GChat) Review {
	return &reviewUsecase{
		historyRepo: historyRepo,
		gchatUc:     gchatUc,
	}
}

// SubmitReviewRequest submits a new review request to tech lead
func (u *reviewUsecase) SubmitReviewRequest(ctx context.Context, req *ReviewRequest, withSend bool) (*ReviewHistoryEntry, error) {
	cfg := config.GetConfig()

	// Validate webhook URL only if sending
	if withSend && cfg.GChatReviewWebhookURL == "" {
		return nil, fmt.Errorf("GChat review webhook URL is not configured")
	}

	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("generate ID: %w", err)
	}

	// Create history entry
	now := time.Now()
	entry := &entity.ReviewHistoryEntry{
		ID:                id,
		Title:             req.Title,
		Description:       req.Description,
		Priority:          req.Priority,
		ReviewLinks:       req.ReviewLinks,
		JiraLinks:         req.JiraLinks,
		SubmittedBy:       cfg.UserName,
		SubmittedByEmail:  cfg.UserEmail,
		SubmittedAt:       now,
		SubmittedToCollab: false,
	}

	// If not sending, return preview only
	if !withSend {
		return entry, nil
	}

	// Save to repository
	if err := u.historyRepo.Save(ctx, entry); err != nil {
		return nil, fmt.Errorf("save history: %w", err)
	}

	// Send to GChat (Tech Lead)
	message := formatReviewRequestMessage(entry)
	if err := u.gchatUc.SendMessage(ctx, cfg.GChatReviewWebhookURL, message); err != nil {
		return nil, fmt.Errorf("send to GChat: %w", err)
	}

	return entry, nil
}

// FormatReviewRequestMessage formats review request for preview/sending
func (u *reviewUsecase) FormatReviewRequestMessage(entry *ReviewHistoryEntry) string {
	return formatReviewRequestMessage(entry)
}

// GetHistories retrieves review histories with optional filter
func (u *reviewUsecase) GetHistories(ctx context.Context, filter HistoryFilter) ([]*ReviewHistoryEntry, error) {
	var entries []*entity.ReviewHistoryEntry
	var err error

	switch filter {
	case HistoryFilterPending:
		entries, err = u.historyRepo.FindByCollabStatus(ctx, false)
	case HistoryFilterCompleted:
		entries, err = u.historyRepo.FindByCollabStatus(ctx, true)
	default:
		entries, err = u.historyRepo.FindAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("get histories: %w", err)
	}

	return entries, nil
}

// GetHistoryByID retrieves a specific history by ID
func (u *reviewUsecase) GetHistoryByID(ctx context.Context, id string) (*ReviewHistoryEntry, error) {
	entry, err := u.historyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get history by ID: %w", err)
	}

	return entry, nil
}

// SubmitToCollaboration forwards a review request to collaboration channel (head architect)
func (u *reviewUsecase) SubmitToCollaboration(ctx context.Context, historyID string) error {
	cfg := config.GetConfig()

	// Validate webhook URL
	if cfg.GChatCollabWebhookURL == "" {
		return fmt.Errorf("GChat collaboration webhook URL is not configured")
	}

	// Get history entry
	entry, err := u.historyRepo.FindByID(ctx, historyID)
	if err != nil {
		return fmt.Errorf("get history: %w", err)
	}

	// Send to GChat (Head Architect)
	message := formatCollaborationMessage(entry)
	if err := u.gchatUc.SendMessage(ctx, cfg.GChatCollabWebhookURL, message); err != nil {
		return fmt.Errorf("send to GChat: %w", err)
	}

	// Update history entry
	now := time.Now()
	entry.SubmittedToCollab = true
	entry.SubmittedToCollabAt = &now
	entry.SubmittedToCollabBy = cfg.UserName

	if err := u.historyRepo.Update(ctx, entry); err != nil {
		return fmt.Errorf("update history: %w", err)
	}

	return nil
}

// SendToGChat sends a message to Google Chat webhook
func (u *reviewUsecase) SendToGChat(ctx context.Context, webhookURL string, message string) error {
	return u.gchatUc.SendMessage(ctx, webhookURL, message)
}

// Helper functions

func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func formatReviewRequestMessage(entry *entity.ReviewHistoryEntry) string {
	msg := fmt.Sprintf("🔍 *New Review Request*\n\n")
	msg += fmt.Sprintf("*Title:* %s\n", entry.Title)
	msg += fmt.Sprintf("*Priority:* %s\n", entry.Priority)
	msg += fmt.Sprintf("*Submitted by:* %s (%s)\n", entry.SubmittedBy, entry.SubmittedByEmail)
	msg += fmt.Sprintf("*Submitted at:* %s\n\n", entry.SubmittedAt.Format("2006-01-02 15:04:05"))

	if entry.Description != "" {
		msg += fmt.Sprintf("*Description:*\n%s\n\n", entry.Description)
	}

	if len(entry.ReviewLinks) > 0 {
		msg += "*Review Links:*\n"
		for _, link := range entry.ReviewLinks {
			msg += fmt.Sprintf("• %s\n", link)
		}
		msg += "\n"
	}

	if len(entry.JiraLinks) > 0 {
		msg += "*Jira Links:*\n"
		for _, link := range entry.JiraLinks {
			msg += fmt.Sprintf("• %s\n", link)
		}
		msg += "\n"
	}

	msg += fmt.Sprintf("*Request ID:* `%s`\n", entry.ID)

	return msg
}

func formatCollaborationMessage(entry *entity.ReviewHistoryEntry) string {
	msg := fmt.Sprintf("🚀 *Review Request*\n\n")
	msg += fmt.Sprintf("*Title:* %s\n", entry.Title)
	msg += fmt.Sprintf("*Priority:* %s\n", entry.Priority)
	msg += fmt.Sprintf("*Originally submitted by:* %s (%s)\n", entry.SubmittedBy, entry.SubmittedByEmail)
	msg += fmt.Sprintf("*Tech Lead Approved:* ✅\n")
	msg += fmt.Sprintf("*Forwarded at:* %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	if entry.Description != "" {
		msg += fmt.Sprintf("*Description:*\n%s\n\n", entry.Description)
	}

	if len(entry.ReviewLinks) > 0 {
		msg += "*Review Links:*\n"
		for _, link := range entry.ReviewLinks {
			msg += fmt.Sprintf("• %s\n", link)
		}
		msg += "\n"
	}

	if len(entry.JiraLinks) > 0 {
		msg += "*Jira Links:*\n"
		for _, link := range entry.JiraLinks {
			msg += fmt.Sprintf("• %s\n", link)
		}
		msg += "\n"
	}

	msg += fmt.Sprintf("*Request ID:* `%s`\n", entry.ID)
	msg += "\n_Please review and approve._"

	return msg
}
