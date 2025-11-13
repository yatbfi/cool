package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/yatbfi/cool/internal/domain/entity"
	domainRepo "github.com/yatbfi/cool/internal/domain/repository"
)

// reviewHistoryRepository implements ReviewHistoryRepository interface
type reviewHistoryRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewReviewHistoryRepository creates a new review history repository
func NewReviewHistoryRepository() (domainRepo.ReviewHistoryRepository, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home dir: %w", err)
	}

	configDir := filepath.Join(home, ".cool-cli")
	filePath := filepath.Join(configDir, "review_histories.json")

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	repo := &reviewHistoryRepository{
		filePath: filePath,
	}

	// Initialize file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := repo.writeHistories([]*entity.ReviewHistoryEntry{}); err != nil {
			return nil, fmt.Errorf("initialize history file: %w", err)
		}
	}

	return repo, nil
}

// Save saves a new review history entry
func (r *reviewHistoryRepository) Save(_ context.Context, entry *entity.ReviewHistoryEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	histories, err := r.readHistories()
	if err != nil {
		return err
	}

	// Check for duplicate ID
	for _, h := range histories {
		if h.ID == entry.ID {
			return fmt.Errorf("entry with ID %s already exists", entry.ID)
		}
	}

	histories = append(histories, entry)
	return r.writeHistories(histories)
}

// Update updates an existing review history entry
func (r *reviewHistoryRepository) Update(_ context.Context, entry *entity.ReviewHistoryEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	histories, err := r.readHistories()
	if err != nil {
		return err
	}

	found := false
	for i, h := range histories {
		if h.ID == entry.ID {
			histories[i] = entry
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("entry with ID %s not found", entry.ID)
	}

	return r.writeHistories(histories)
}

// FindByID retrieves a review history entry by ID
func (r *reviewHistoryRepository) FindByID(_ context.Context, id string) (*entity.ReviewHistoryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	histories, err := r.readHistories()
	if err != nil {
		return nil, err
	}

	for _, h := range histories {
		if h.ID == id {
			return h, nil
		}
	}

	return nil, fmt.Errorf("entry with ID %s not found", id)
}

// FindAll retrieves all review history entries
func (r *reviewHistoryRepository) FindAll(_ context.Context) ([]*entity.ReviewHistoryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.readHistories()
}

// FindByCollabStatus retrieves review history entries filtered by collaboration status
func (r *reviewHistoryRepository) FindByCollabStatus(_ context.Context, submittedToCollab bool) ([]*entity.ReviewHistoryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	histories, err := r.readHistories()
	if err != nil {
		return nil, err
	}

	var filtered []*entity.ReviewHistoryEntry
	for _, h := range histories {
		if h.SubmittedToCollab == submittedToCollab {
			filtered = append(filtered, h)
		}
	}

	return filtered, nil
}

// Delete deletes a review history entry by ID
func (r *reviewHistoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	histories, err := r.readHistories()
	if err != nil {
		return err
	}

	found := false
	newHistories := make([]*entity.ReviewHistoryEntry, 0, len(histories))
	for _, h := range histories {
		if h.ID != id {
			newHistories = append(newHistories, h)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("entry with ID %s not found", id)
	}

	return r.writeHistories(newHistories)
}

// readHistories reads all histories from JSON file
func (r *reviewHistoryRepository) readHistories() ([]*entity.ReviewHistoryEntry, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("read history file: %w", err)
	}

	var histories []*entity.ReviewHistoryEntry
	if err := json.Unmarshal(data, &histories); err != nil {
		return nil, fmt.Errorf("unmarshal histories: %w", err)
	}

	return histories, nil
}

// writeHistories writes all histories to JSON file
func (r *reviewHistoryRepository) writeHistories(histories []*entity.ReviewHistoryEntry) error {
	data, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal histories: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0o644); err != nil {
		return fmt.Errorf("write history file: %w", err)
	}

	return nil
}
