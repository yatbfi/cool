# AI Prompting Guide for Review History Feature

## Quick Reference for AI Assistants

This document provides quick context for AI assistants working on the review history feature.

---

## Feature Summary

**Purpose**: Track code review requests and enable forwarding to collaboration channel

**Flow**: Developer → Tech Lead (review) → Head Architect (collaboration)

**Commands**:
- `cool review request` - Submit review to Tech Lead
- `cool review histories` - View all history
- `cool review submit-collab` - Forward to Head Architect

---

## Architecture Map

```
Domain Layer (usecase/):
├── review.go
│   ├── ReviewHistoryEntry (model)
│   ├── ReviewHistoryRepository (port/interface)
│   ├── Review (port/interface)
│   └── ReviewUsecase (implementation)

Infrastructure Layer (infrastructure/repository/):
└── review_history.go
    └── FileReviewHistoryRepository (adapter)

Presentation Layer (cmd/):
├── review_request.go (existing, updated)
├── review_histories.go (new)
└── review_submit_collab.go (new)

Config:
└── ~/.cool-cli/review_histories.json (data file)
```

---

## Key Files

### Domain Layer
**File**: `internal/domain/usecase/review.go`
- Contains business logic
- Defines interfaces (ports)
- Independent of infrastructure

**Key Types**:
```go
type ReviewHistoryEntry struct { ... }
type ReviewHistoryRepository interface { ... }
type Review interface { ... }
type ReviewUsecase struct { ... }
```

### Infrastructure Layer
**File**: `internal/infrastructure/repository/review_history.go`
- Implements ReviewHistoryRepository
- Uses JSON file storage
- Located: `~/.cool-cli/review_histories.json`

### Presentation Layer
**Files**: `cmd/review_histories.go`, `cmd/review_submit_collab.go`
- CLI command handlers
- User interaction
- Delegates to usecase

### Main
**File**: `main.go`
- Dependency injection
- Command registration

---

## Common Prompting Scenarios

### Scenario 1: Modify History Display

**Prompt**: "Change the history table to show XXX"

**Context needed**:
- File: `cmd/review_histories.go`
- Function: `displayHistoriesTable()`
- Uses: `text/tabwriter` for formatting

### Scenario 2: Add New Field to History

**Steps**:
1. Add field to `ReviewHistoryEntry` in `usecase/review.go`
2. Update JSON tags
3. Update display in `cmd/review_histories.go`
4. Update preview in `cmd/review_submit_collab.go`
5. Update message builder if needed

### Scenario 3: Change Storage Format

**Context needed**:
- Current: JSON file in `~/.cool-cli/`
- Repository interface: Already abstracted
- Only need to implement new adapter

**Example prompt**: "Change storage from JSON to SQLite"
**What to modify**: Only `internal/infrastructure/repository/review_history.go`

### Scenario 4: Add New Filter

**Steps**:
1. Add filter to `HistoryFilter` enum in `usecase/review.go`
2. Add repository method if needed
3. Update `GetHistories()` switch statement
4. Add flag to `cmd/review_histories.go`

### Scenario 5: Modify Message Format

**Context needed**:
- Review message: `buildReviewMessage()` in `usecase/review.go`
- Collab message: `buildCollaborationMessage()` in `usecase/review.go`
- Format: Google Chat markdown

---

## File Locations Quick Reference

```
Domain Logic:
  internal/domain/usecase/review.go

Infrastructure:
  internal/infrastructure/repository/review_history.go

Commands:
  cmd/review_request.go      (updated)
  cmd/review_histories.go    (new)
  cmd/review_submit_collab.go (new)
  cmd/base.go                (updated - validation)

Main:
  main.go                    (updated - DI)

Documentation:
  docs/REVIEW_HISTORY_SPEC.md
  docs/REVIEW_HISTORY_GUIDE.md
  docs/REVIEW_HISTORY_TESTING.md

Data:
  ~/.cool-cli/review_histories.json (runtime)
```

---

## Design Principles

### Hexagonal Architecture
- **Domain**: Business logic, isolated
- **Ports**: Interfaces (Repository, GChat)
- **Adapters**: Implementations (FileRepository, Commands)

### Dependency Injection
All dependencies injected in `main.go`:
```go
historyRepo := repository.NewFileReviewHistoryRepository()
reviewUc := usecase.NewReviewUsecase(cfg, gChatUc, historyRepo)
```

### Single Responsibility
- Commands: User interaction only
- usecase: Business logic only
- Repository: Data persistence only

---

## Common Modifications

### Add New Command

1. Create `cmd/review_xxx.go`
2. Implement command struct with `baseCmd`
3. Inject usecase in constructor
4. Register in `main.go`:
   ```go
   cmd.NewReviewXxxCmd(reviewUc).Reg(reviewCmd)
   ```

### Add New Validation

In `cmd/base.go`, update `validateCommandSpecificConfig()`:
```go
case "new-command":
    // validation logic
```

### Modify History Schema

1. Update `ReviewHistoryEntry` struct
2. Add JSON tags
3. Consider migration if breaking change
4. Update displays and messages

### Add New Repository Method

1. Add to interface in `usecase/review.go`
2. Implement in `repository/review_history.go`
3. Use in usecase method

---

## Testing Prompts

### Unit Test Request
"Write unit tests for the review history repository"

**What to test**:
- Save, FindAll, FindByID, FindPending, FindCompleted, Update
- Edge cases: empty file, corrupted JSON, missing ID

### Integration Test Request
"Write integration test for submit-collab flow"

**What to test**:
- Submit review → Save history → Submit to collab → Update status

---

## Data Flow

### Submit Review Request
```
User Input (cmd)
    ↓
ReviewUsecase.RequestReview()
    ↓
├─→ Send to GChat (lora_code_review)
└─→ Save to History (historyRepo.Save())
```

### Submit to Collaboration
```
User Selection (cmd)
    ↓
ReviewUsecase.SubmitToCollaboration()
    ↓
├─→ Find by ID (historyRepo.FindByID())
├─→ Validate not already submitted
├─→ Send to GChat (lora_collaboration)
└─→ Update history (historyRepo.Update())
```

### View Histories
```
User Request (cmd)
    ↓
ReviewUsecase.GetHistories(filter)
    ↓
historyRepo.FindAll/FindPending/FindCompleted()
    ↓
Display Table (cmd)
```

---

## Error Handling Patterns

### Domain Errors
Return error with context:
```go
return fmt.Errorf("failed to send review message: %w", err)
```

### User-Facing Errors
Show helpful message:
```go
fmt.Println("⚠️  GChat review webhook URL is not configured.")
fmt.Println("Please run the setup command to configure it:")
return fmt.Errorf("missing GChat review webhook URL")
```

### Non-Critical Errors
Warn but continue:
```go
if err := u.historyRepo.Save(ctx, entry); err != nil {
    fmt.Printf("⚠️  Warning: Failed to save history: %v\n", err)
}
```

---

## Message Format Reference

### Google Chat Markdown
```
*Bold*
_Italic_
<url|text>
```

### Review Message Template
```
*[P1] Title*

_Description_

*Pull Requests:*
• Service → <url|#123>

*Jira Links:*
• <url|TICKET-123>
```

### Collaboration Message Template
Same as review + these fields:
```
*Originally Submitted:* 2025-11-13 10:30
*Submitted by:* Name (email)
```

---

## Common Issues & Solutions

### "Review not found"
- Check ID format (support both full and short ID)
- Implement fuzzy matching in `FindByID()`

### "Already submitted"
- Check `SubmittedToCollab` flag before submitting
- Show clear error with timestamp

### "Table formatting broken"
- Ensure using `text/tabwriter`
- Check terminal width
- Truncate long strings

### "History file corrupted"
- Backup file before writing
- Handle JSON parse errors gracefully
- Provide recovery instructions

---

## Extension Points

### Easy to Add
- New filter types (by priority, by date, etc.)
- Export functionality (CSV, Markdown)
- Statistics view
- Notes/comments on reviews

### Requires More Work
- Database storage (need new adapter)
- Multi-user support (need user ID tracking)
- Webhooks for notifications
- Web UI

---

## Validation Rules

### Review Request
- Must have title
- Must have at least 1 PR link
- Must have description
- Must have priority
- Must have review webhook configured

### Submit to Collab
- Must have collab webhook configured
- Review must exist
- Review must not be already submitted
- User must confirm

---

## Configuration Keys

In `~/.cool-cli/config.json`:
```json
{
  "gchat_review_webhook_url": "...",      // Required for review request
  "gchat_collab_webhook_url": "...",      // Required for submit-collab
  "user_name": "...",                     // Required for both
  "user_email": "..."                     // Required for both
}
```

---

## Useful Grep Patterns

Find all review-related commands:
```bash
grep -r "NewReview" cmd/
```

Find all history repository usages:
```bash
grep -r "historyRepo" internal/
```

Find all GChat sends:
```bash
grep -r "SendToChannel" internal/
```

---

## Quick Modification Recipes

### Add Priority Filter
```go
// 1. Add flag in cmd/review_histories.go
flags.StringVar(&c.priority, "priority", "", "Filter by priority (P0,P1,P2,P3)")

// 2. Add repository method
func (r *FileReviewHistoryRepository) FindByPriority(ctx context.Context, priority string) ([]*usecase.ReviewHistoryEntry, error) { ... }

// 3. Update usecase if needed
```

### Add Export Command
```go
// 1. Create cmd/review_export.go
// 2. Implement CSV/Markdown writer
// 3. Use GetHistories() to fetch data
// 4. Register in main.go
```

### Add Search
```go
// 1. Add flag in cmd/review_histories.go
flags.StringVar(&c.search, "search", "", "Search in title/description")

// 2. Filter results after GetHistories()
// Or add FindBySearch() method
```

---

## Performance Considerations

- History file loads into memory (fast for < 1000 entries)
- Sorting done in-memory (negligible for normal use)
- No pagination needed for CLI use case
- File I/O is bottleneck (use defer for writes)

---

## Security Considerations

- History file: 0600 permissions (owner only)
- No sensitive data in history (only URLs and titles)
- Webhook URLs in separate config file
- No logging of webhook URLs

---

## Future-Proofing

The architecture supports:
- ✅ Changing storage (just implement new adapter)
- ✅ Adding fields to history (JSON is flexible)
- ✅ Adding new commands (pluggable)
- ✅ Adding new filters (enum-based)
- ✅ Adding webhooks (interface-based)

Migration strategy:
1. Add version field to history entries
2. Implement migration functions
3. Run on load if version mismatch

---

## For Code Review

Check these aspects:
- [ ] Error handling complete
- [ ] User messages clear and helpful
- [ ] Domain logic in usecase, not cmd
- [ ] Repository interface clean
- [ ] No business logic in repository
- [ ] Commands thin (just user interaction)
- [ ] Tests written
- [ ] Documentation updated

---

## Example Prompts for AI

### Good Prompts
✅ "Add a search feature to filter histories by keyword"
✅ "Change the history table to sort by priority instead of time"
✅ "Add a --json flag to output histories as JSON"
✅ "Implement export to CSV functionality"

### Prompts Needing Context
⚠️ "Make the table look better" (specify what to improve)
⚠️ "Add a new field" (specify which field and where to display)
⚠️ "Fix the bug" (specify which bug)

### Complex Prompts (Break Down)
🔴 "Add database support, user authentication, and web UI"
✅ Better: Start with "Replace file storage with SQLite"

---

## Cheat Sheet

**View all code**:
```bash
cat internal/domain/usecase/review.go
cat internal/infrastructure/repository/review_history.go
cat cmd/review_histories.go
cat cmd/review_submit_collab.go
```

**Test commands**:
```bash
cool review request --dry-run
cool review histories
cool review histories --pending
cool review submit-collab --help
```

**Check data**:
```bash
cat ~/.cool-cli/review_histories.json | jq .
```

**Build and test**:
```bash
go build -o bin/cool && ./bin/cool review histories
```

---

This guide should help AI assistants quickly understand the codebase and make appropriate modifications. For detailed specifications, refer to `REVIEW_HISTORY_SPEC.md`.

