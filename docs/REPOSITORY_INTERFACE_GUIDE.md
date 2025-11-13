# Repository Interface & Implementation Guide

## Overview

The Review History Repository follows **Interface-Based Design** with clear separation between interface (port) and implementation (adapter), adhering to hexagonal architecture principles.

---

## Architecture

### Hexagonal Architecture (Ports & Adapters)

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer (Core)                       │
│                                                               │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Interface (Port)                                   │    │
│  │  usecase/ReviewHistoryRepository                   │    │
│  │  - Save()                                           │    │
│  │  - FindAll()                                        │    │
│  │  - FindByID()                                       │    │
│  │  - FindPending()                                    │    │
│  │  - FindCompleted()                                  │    │
│  │  - Update()                                         │    │
│  └────────────────────────────────────────────────────┘    │
│                                                               │
└─────────────────────────┬───────────────────────────────────┘
                          │ depends on
                          │ (dependency inversion)
┌─────────────────────────▼───────────────────────────────────┐
│               Infrastructure Layer (Adapters)                │
│                                                               │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Implementation (Adapter)                           │    │
│  │  repository/FileReviewHistoryRepository             │    │
│  │  - Implements all interface methods                 │    │
│  │  - Uses JSON file storage                           │    │
│  │  - Located: ~/.cool-cli/review_histories.json      │    │
│  └────────────────────────────────────────────────────┘    │
│                                                               │
│  Future Implementations:                                     │
│  - DatabaseReviewHistoryRepository (SQLite, Postgres)       │
│  - CloudReviewHistoryRepository (S3, GCS)                   │
│  - InMemoryReviewHistoryRepository (testing)                │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Interface Definition

**Location:** `internal/domain/usecase/review.go`

### Why in `usecase` package?

The interface is defined in the domain layer (usecase) because:
1. **Domain owns the interface** - The business logic defines what it needs
2. **Dependency Inversion** - Infrastructure depends on domain, not vice versa
3. **Avoid Cyclic Imports** - Keeps domain models and interfaces in one place
4. **Hexagonal Architecture** - Ports are defined by the core, not by adapters

### Interface Contract

```go
type ReviewHistoryRepository interface {
    // Save persists a new review history entry
    Save(ctx context.Context, entry *ReviewHistoryEntry) error

    // FindAll retrieves all review history entries
    FindAll(ctx context.Context) ([]*ReviewHistoryEntry, error)

    // FindByID retrieves a specific review by its ID
    // Must support both full ID and short ID (first 8 characters)
    FindByID(ctx context.Context, id string) (*ReviewHistoryEntry, error)

    // FindPending retrieves reviews not submitted to collaboration
    FindPending(ctx context.Context) ([]*ReviewHistoryEntry, error)

    // FindCompleted retrieves reviews submitted to collaboration
    FindCompleted(ctx context.Context) ([]*ReviewHistoryEntry, error)

    // Update updates an existing review history entry
    Update(ctx context.Context, entry *ReviewHistoryEntry) error
}
```

---

## Current Implementation

**Location:** `internal/infrastructure/repository/review_history.go`

### FileReviewHistoryRepository

```go
type FileReviewHistoryRepository struct {
    filePath string
}

// Compile-time interface verification
var _ usecase.ReviewHistoryRepository = (*FileReviewHistoryRepository)(nil)
```

**Features:**
- ✅ JSON file-based storage
- ✅ Located at `~/.cool-cli/review_histories.json`
- ✅ Secure permissions (0600 - owner only)
- ✅ Pretty-printed JSON for readability
- ✅ Automatic sorting (newest first)
- ✅ Support for short IDs (8 characters)
- ✅ Graceful error handling
- ✅ Automatic backup on corruption

**Storage Format:**
```json
[
  {
    "id": "1731488400-123456789",
    "title": "[P1] Feature",
    "description": "Description",
    "review_links": [...],
    "jira_links": [...],
    "priority": "P1",
    "submitted_at": "2025-11-13T10:30:00+07:00",
    "submitted_to_collab": false,
    "submitted_by": "User (email)"
  }
]
```

---

## Design Benefits

### 1. Testability

Easy to mock for testing:

```go
type MockReviewHistoryRepository struct {
    SaveFunc       func(ctx context.Context, entry *usecase.ReviewHistoryEntry) error
    FindAllFunc    func(ctx context.Context) ([]*usecase.ReviewHistoryEntry, error)
    // ... other methods
}

func (m *MockReviewHistoryRepository) Save(ctx context.Context, entry *usecase.ReviewHistoryEntry) error {
    return m.SaveFunc(ctx, entry)
}
```

### 2. Swappable Implementations

Easy to switch storage without changing business logic:

```go
// Development: Use file storage
historyRepo := repository.NewFileReviewHistoryRepository()

// Production: Use database
historyRepo := repository.NewDatabaseReviewHistoryRepository(db)

// Testing: Use in-memory
historyRepo := repository.NewInMemoryReviewHistoryRepository()

// Same usecase, different storage!
reviewUc := usecase.NewReviewUsecase(cfg, gChat, historyRepo)
```

### 3. Dependency Inversion

Business logic depends on abstraction (interface), not concrete implementation:

```go
// ReviewUsecase depends on interface
type ReviewUsecase struct {
    historyRepo usecase.ReviewHistoryRepository // Interface, not concrete type!
}
```

### 4. Clean Architecture

Clear separation of concerns:
- **Domain** - Defines what it needs (interface)
- **Infrastructure** - Provides what domain needs (implementation)
- **Presentation** - Uses domain, doesn't care about implementation

---

## Creating New Implementations

### Example: Database Implementation

**Step 1:** Create new file

```go
// internal/infrastructure/repository/review_history_db.go
package repository

import (
    "context"
    "database/sql"
    "github.com/yatbfi/cool/internal/domain/usecase"
)

type DatabaseReviewHistoryRepository struct {
    db *sql.DB
}

// Verify interface implementation
var _ usecase.ReviewHistoryRepository = (*DatabaseReviewHistoryRepository)(nil)

func NewDatabaseReviewHistoryRepository(db *sql.DB) *DatabaseReviewHistoryRepository {
    return &DatabaseReviewHistoryRepository{db: db}
}

func (r *DatabaseReviewHistoryRepository) Save(ctx context.Context, entry *usecase.ReviewHistoryEntry) error {
    query := `INSERT INTO review_histories (...) VALUES (...)`
    _, err := r.db.ExecContext(ctx, query, ...)
    return err
}

// Implement all other interface methods...
```

**Step 2:** Update dependency injection in main.go

```go
// Switch implementation here
db, _ := sql.Open("postgres", connectionString)
historyRepo := repository.NewDatabaseReviewHistoryRepository(db)

// Business logic stays the same!
reviewUc := usecase.NewReviewUsecase(cfg, gChat, historyRepo)
```

### Example: In-Memory Implementation (for testing)

```go
// internal/infrastructure/repository/review_history_memory.go
package repository

import (
    "context"
    "fmt"
    "sync"
    "github.com/yatbfi/cool/internal/domain/usecase"
)

type InMemoryReviewHistoryRepository struct {
    mu      sync.RWMutex
    entries map[string]*usecase.ReviewHistoryEntry
}

var _ usecase.ReviewHistoryRepository = (*InMemoryReviewHistoryRepository)(nil)

func NewInMemoryReviewHistoryRepository() *InMemoryReviewHistoryRepository {
    return &InMemoryReviewHistoryRepository{
        entries: make(map[string]*usecase.ReviewHistoryEntry),
    }
}

func (r *InMemoryReviewHistoryRepository) Save(ctx context.Context, entry *usecase.ReviewHistoryEntry) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.entries[entry.ID] = entry
    return nil
}

// Implement other methods...
```

---

## Dependency Injection

### Current Setup (main.go)

```go
func initCommands(cfg *config.Config) {
    // 1. Create repository (infrastructure layer)
    historyRepo := infraRepo.NewFileReviewHistoryRepository()

    // 2. Inject into usecase (domain layer)
    gChatUc := usecase.NewGChatUsecase(cfg)
    reviewUc := usecase.NewReviewUsecase(cfg, gChatUc, historyRepo)

    // 3. Inject into commands (presentation layer)
    reviewCmd := cmd.NewReviewCmd().Reg(rootCmd)
    cmd.NewReviewRequestCmd(reviewUc).Reg(reviewCmd)
    cmd.NewReviewHistoriesCmd(reviewUc).Reg(reviewCmd)
    cmd.NewReviewSubmitCollabCmd(reviewUc).Reg(reviewCmd)
}
```

### Flow of Dependencies

```
main.go (composition root)
    │
    ├─ Creates: FileReviewHistoryRepository
    │           (infrastructure/adapter)
    │
    ├─ Injects into: ReviewUsecase
    │                (domain/usecase)
    │
    └─ Injects into: Commands
                     (presentation/cmd)
```

---

## Interface Verification

### Compile-Time Check

The line:
```go
var _ usecase.ReviewHistoryRepository = (*FileReviewHistoryRepository)(nil)
```

Ensures that `FileReviewHistoryRepository` implements all methods of the interface. If any method is missing or has wrong signature, compilation will fail.

**Benefits:**
- ✅ Catch missing methods at compile time
- ✅ Catch wrong signatures at compile time
- ✅ Self-documenting code
- ✅ Refactoring safety

---

## Best Practices

### 1. Keep Interface Minimal

Only define methods that are actually needed. Don't add "just in case" methods.

✅ Good:
```go
type ReviewHistoryRepository interface {
    Save(ctx context.Context, entry *ReviewHistoryEntry) error
    FindByID(ctx context.Context, id string) (*ReviewHistoryEntry, error)
}
```

❌ Bad:
```go
type ReviewHistoryRepository interface {
    Save(...)
    FindByID(...)
    FindByTitle(...)           // Not needed yet
    FindByDate(...)            // Not needed yet
    DeleteByID(...)            // Not needed yet
    BulkInsert(...)           // Not needed yet
}
```

### 2. Use Context

Always accept `context.Context` as first parameter for cancellation and timeout support.

```go
func (r *FileReviewHistoryRepository) Save(ctx context.Context, entry *ReviewHistoryEntry) error {
    // Can check ctx.Done() for cancellation
}
```

### 3. Return Errors

Always return errors for proper error handling:

```go
// Good
func (r *Repo) Save(...) error

// Bad
func (r *Repo) Save(...) bool
```

### 4. Consistent Naming

- `Save` - Create new entry
- `Find*` - Query entries
- `Update` - Modify existing entry
- `Delete` - Remove entry (if needed)

### 5. Document Interface

Add comments to interface methods explaining:
- What the method does
- What parameters mean
- What errors might be returned
- Any special behavior

---

## Testing Strategy

### Unit Tests for Implementation

Test each repository implementation independently:

```go
func TestFileReviewHistoryRepository_Save(t *testing.T) {
    // Setup temp directory
    repo := NewFileReviewHistoryRepository()
    
    // Test save
    entry := &usecase.ReviewHistoryEntry{...}
    err := repo.Save(context.Background(), entry)
    
    assert.NoError(t, err)
}
```

### Integration Tests

Test with real storage:

```go
func TestReviewWorkflow_WithFileStorage(t *testing.T) {
    // Use real file repository
    repo := repository.NewFileReviewHistoryRepository()
    usecase := usecase.NewReviewUsecase(cfg, gChat, repo)
    
    // Test full workflow
    // ...
}
```

### Mock for Usecase Tests

Mock repository for testing business logic:

```go
func TestReviewUsecase_SubmitToCollaboration(t *testing.T) {
    // Mock repository
    mockRepo := &MockReviewHistoryRepository{
        FindByIDFunc: func(ctx context.Context, id string) (*usecase.ReviewHistoryEntry, error) {
            return &usecase.ReviewHistoryEntry{...}, nil
        },
    }
    
    usecase := usecase.NewReviewUsecase(cfg, gChat, mockRepo)
    
    // Test business logic without real storage
    // ...
}
```

---

## Summary

### Key Points

1. **Interface in Domain** - `usecase.ReviewHistoryRepository`
2. **Implementation in Infrastructure** - `repository.FileReviewHistoryRepository`
3. **Dependency Injection** - In `main.go` (composition root)
4. **Interface Verification** - Compile-time check with `var _`
5. **Swappable** - Easy to change implementation
6. **Testable** - Easy to mock and test

### Benefits Achieved

✅ **Separation of Concerns** - Clear boundaries  
✅ **Dependency Inversion** - Domain independent of infrastructure  
✅ **Testability** - Easy to mock and test  
✅ **Flexibility** - Easy to swap implementations  
✅ **Maintainability** - Changes isolated to appropriate layer  

### Current State

- ✅ Interface defined in domain
- ✅ File-based implementation working
- ✅ Dependency injection configured
- ✅ Compile-time verification in place
- ✅ Ready for additional implementations

---

**This design makes the codebase flexible, testable, and maintainable!** 🚀

