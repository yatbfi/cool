package repository

import (
	"context"

	"github.com/yatbfi/cool/internal/domain/entity"
)

// ReviewHistoryRepository defines the interface for managing review history
type ReviewHistoryRepository interface {
	// Save saves a new review history entry
	Save(ctx context.Context, entry *entity.ReviewHistoryEntry) error

	// Update updates an existing review history entry
	Update(ctx context.Context, entry *entity.ReviewHistoryEntry) error

	// FindByID retrieves a review history entry by ID
	FindByID(ctx context.Context, id string) (*entity.ReviewHistoryEntry, error)

	// FindAll retrieves all review history entries
	FindAll(ctx context.Context) ([]*entity.ReviewHistoryEntry, error)

	// FindByCollabStatus retrieves review history entries filtered by collaboration status
	FindByCollabStatus(ctx context.Context, submittedToCollab bool) ([]*entity.ReviewHistoryEntry, error)

	// Delete deletes a review history entry by ID
	Delete(ctx context.Context, id string) error
}
