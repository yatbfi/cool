# Cool CLI - AI Coding Guide

## Architecture Overview

This is a **Go CLI tool** following **Clean Architecture (Hexagonal)** principles. The codebase is strictly layered:

```
cmd/                    # Presentation: CLI commands (Cobra)
internal/domain/        # Business logic (no external dependencies)
├── entity/            # Domain models
├── repository/        # Port interfaces
└── usecase/           # Business logic orchestration
internal/infrastructure/  # Adapters: implementations of ports
└── repository/        # Concrete repository implementations
config/                # Configuration management (~/.cool-cli/)
```

**Critical Rule**: Dependencies flow inward only. Domain layer NEVER imports from infrastructure or cmd layers.

## Key Patterns

### 1. Dependency Injection in root.go
All dependencies are wired in `cmd/root.go`:
```go
historyRepo := infraRepo.NewReviewHistoryRepository()
gchatUc := usecase.NewGChatUsecase()
reviewUc := usecase.NewReviewUsecase(historyRepo, gchatUc)
```

When adding features: define interface in `domain/usecase/` or `domain/repository/`, implement in `infrastructure/repository/`, inject in `root.go`.

### 2. Command Structure Pattern
All commands embed `*baseCmd` which provides:
- Pre-run validation (environment, user setup, command-specific config)
- Skip validation for `setup` and `update` commands via `shouldSkipValidation()`
- See `cmd/base.go` for validation logic

Example: `review request` requires `GChatReviewWebhookURL`, `submit-collab` requires `GChatCollabWebhookURL`.

### 3. Configuration Management
- Location: `~/.cool-cli/config.json` (user config), `~/.cool-cli/review_histories.json` (review data)
- Access: `config.GetConfig()` returns cached singleton
- Update: `config.SaveLocalConfig()` persists changes
- Fields: `UserName`, `UserEmail`, `GChatReviewWebhookURL`, `GChatCollabWebhookURL`, `PreferredEditor`, `ProjectRoot`

### 4. Repository Interface Pattern
Repository interfaces live in `internal/domain/repository/`, implementations in `internal/infrastructure/repository/`.

Example (`review_history.go`):
- Interface: `domain/repository/review_history.go`
- Implementation: `infrastructure/repository/review_history.go` (JSON file storage)
- Thread-safe with `sync.RWMutex`

## Development Workflows

### Building & Running
```bash
make build          # Build to bin/cool
make install        # Install to $GOBIN
make lint           # Run golangci-lint
make all            # Cross-compile for Linux/macOS/Windows
```

### Hot Reload Development
```bash
cool run -- go run .                    # Hot reload current directory
cool run --path=./cmd -- go run ./cmd  # Watch specific path
cool run --debug -- go run .            # Enable debug logging
```

Hot reload uses polling-based file watcher (500ms interval) - see `internal/infrastructure/repository/file_watcher.go`.

### Testing
No test files exist currently. When adding tests, mock interfaces from `domain/repository/` and `domain/usecase/`.

## Feature Implementation Checklist

Adding a new review-related feature:

1. **Entity**: Define model in `internal/domain/entity/` (if needed)
2. **Repository Interface**: Define port in `internal/domain/repository/`
3. **Repository Implementation**: Implement in `internal/infrastructure/repository/`
4. **Usecase Interface**: Define business logic interface in `internal/domain/usecase/`
5. **Usecase Implementation**: Implement orchestration logic
6. **Command**: Create cmd handler in `cmd/`, embed `*baseCmd`
7. **Wire Dependencies**: Inject in `cmd/root.go` `NewRootCommand()`
8. **Validation**: Add command-specific validation in `cmd/base.go` if needed

## Important Conventions

- **Error Handling**: Return wrapped errors with context: `fmt.Errorf("description: %w", err)`
- **Context**: All repository/usecase methods accept `context.Context` as first param
- **Thread Safety**: Repositories managing shared state use `sync.RWMutex`
- **Editor Integration**: Use `internal/pkg/common/OpenEditorForInput()` for multiline input (respects `PreferredEditor` config)
- **Table Display**: Use `text/tabwriter` for formatted output (see `cmd/review_histories.go`)
- **ID Generation**: Use `crypto/rand` for unique IDs (see `usecase/review.go:generateID()`)

## External Dependencies

- **Cobra**: CLI framework (`github.com/spf13/cobra`)
- **Logrus**: Logging (`github.com/sirupsen/logrus`) - used in hot reload only
- **PromptUI**: Interactive prompts (`github.com/manifoldco/promptui`) - used in setup commands
- **Google Chat**: Webhook-based integration (no SDK) - see `usecase/g_chat.go`

## Common Pitfalls

- Don't bypass `baseCmd` validation - it ensures required config exists before command execution
- Don't import from higher layers (infrastructure → domain is forbidden)
- Don't forget to inject new dependencies in `root.go`
- Review history operations must lock with mutex (see `infrastructure/repository/review_history.go`)
- Editor command splits on spaces - handle `code --wait` style commands correctly (see `common/editor.go`)
