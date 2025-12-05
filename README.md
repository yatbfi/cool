# Cool CLI - Developer Review Workflow Tool 🚀

A powerful command-line tool to streamline the code review workflow from developer to tech lead to head architect, with seamless Google Chat integration and hot reload development tooling.

## ✨ Features

- 📝 **Review Request Submission** - Submit review requests to tech lead with all necessary details
- ✏️ **Multiline Input with Editor** - Use your preferred editor (vim, nano, etc.) for descriptions
- 👀 **Preview & Edit Before Submit** - Preview formatted message and edit any field before sending
- 🎯 **P0-P4 Priority System** - Clear priority levels from Critical to Very Low
- 📊 **History Tracking** - Track all your review requests with comprehensive history
- 🔄 **Collaboration Forwarding** - Forward approved reviews to head architect
- 💬 **Google Chat Integration** - Automatic notifications to review channels
- 🎨 **Beautiful Table Display** - Clean and organized view of your reviews
- ⚙️ **Easy Configuration** - Simple setup wizard for user info, webhooks, editor, and project root
- 🔥 **Hot Reload** - Automatic application restart on file changes during development
- 🏗️ **Clean Architecture** - Domain-driven design with clear separation of concerns
- 🛡️ **Thread-Safe** - Safe concurrent operations with mutex protection

## 🏗️ Architecture

This project follows Clean Architecture (Hexagonal Architecture) principles:

```
cmd/                    # Presentation layer (CLI commands)
├── root.go
├── setup.go
├── review_request.go
├── run.go
└── ...

internal/
├── domain/             # Business logic & interfaces
│   ├── entity/        # Domain entities (models)
│   ├── repository/    # Repository interfaces (ports)
│   └── usecase/       # Business use cases
│
├── infrastructure/     # External implementations (adapters)
│   └── repository/    # Concrete repository implementations
│
└── pkg/               # Shared utilities
    ├── common/        # Common helpers
    ├── logger/        # Logging utilities
    └── table/         # Table display

config/                 # Configuration management
```

**Layer Responsibilities:**
- **Domain Layer**: Pure business logic, no external dependencies
- **Infrastructure Layer**: Concrete implementations of repositories, external APIs
- **Presentation Layer**: CLI commands, user interaction
- **Use Cases**: Orchestrate business logic using domain entities and repositories

## 🚀 Quick Start

### Installation

```bash
# Install from source
go install github.com/yatbfi/cool@latest

# Or build locally
git clone https://github.com/yatbfi/cool.git
cd cool
go build -o bin/cool .
```

### Initial Setup

```bash
# Run full setup wizard
cool setup

# Or configure components separately
cool setup email          # Configure your name and email
cool setup webhook        # Configure Google Chat webhooks
cool setup editor         # Choose preferred text editor
cool setup project-root   # Set workspace project root

# View current configuration
cool config preview
```

### Basic Usage

```bash
# Submit a review request (opens editor for description)
cool review request

# View all review history
cool review histories

# View pending reviews only
cool review histories --pending

# Submit to collaboration (after tech lead approval)
cool review submit-collab <review-id>

# Run with hot reload from project root
cd /path/to/your/project
cool run -- go run ./cmd/http

# Run with .env auto-load
cool run -- go run cmd/http/main.go

# Hot reload with custom patterns
cool run --regex=".*\\.go$,.*\\.html$" -- go run .
```

## 📋 Commands

### Setup Commands

| Command | Description | Interactive |
|---------|-------------|-------------|
| `cool setup` | Run full setup wizard | ✅ |
| `cool setup email` | Configure user name and email | ✅ |
| `cool setup webhook` | Configure Google Chat webhook URLs | ✅ |
| `cool setup editor` | Choose preferred text editor | ✅ |
| `cool setup project-root` | Set workspace project root | ✅ |
| `cool config preview` | Display all configuration settings | ❌ |

### Review Commands

| Command | Description | Notes |
|---------|-------------|-------|
| `cool review request` | Submit new review request to tech lead | Opens editor for description |
| `cool review histories` | Display all review history | Shows all statuses |
| `cool review histories --pending` | Show pending reviews only | Filter by status |
| `cool review histories --completed` | Show completed reviews only | Filter by status |
| `cool review submit-collab <id>` | Submit review to head architect | Requires review ID |
| `cool review submit-collab --list` | List all reviews | Alternative to histories |
| `cool review submit-collab --list --pending` | List pending reviews | Combined filter |

### Hot Reload Commands

| Command | Description | Example |
|---------|-------------|---------|
| `cool run -- [command]` | Run with hot reload | `cool run -- go run .` |
| `cool run -- go run <path>` | Run specific main.go | `cool run -- go run cmd/http/main.go` |
| `cool run --path=<dir> -- [cmd]` | Watch specific directory | `cool run --path=./services -- go run .` |
| `cool run --regex=<patterns> -- [cmd]` | Custom file patterns | `cool run --regex=".*\\.go$,.*\\.html$" -- go run .` |
| `cool run --env=<vars> -- [cmd]` | Set environment variables | `cool run --env="PORT=8080,DEBUG=true" -- go run .` |
| `cool run --debug -- [cmd]` | Enable debug logging | `cool run --debug -- go test` |
| `cool run --quiet -- [cmd]` | Disable logging | `cool run --quiet -- go run .` |

**Note**: Use `--` separator between flags and command for clarity.

### Other Commands

| Command | Description |
|---------|-------------|
| `cool update` | Update Cool CLI to latest version |
| `cool completion` | Generate shell completion scripts |

## 🎯 Workflow

### Review Request Workflow

```
Developer                     Tech Lead              Head Architect
    │                             │                        │
    │ 1. cool review request      │                        │
    ├────────────────────────────>│                        │
    │    (Opens editor)            │                        │
    │    (Preview & Edit)          │                        │
    │    (Google Chat notify)      │                        │
    │                             │                        │
    │                             │ 2. Review & Approve    │
    │                             │                        │
    │ 3. cool review submit-collab│                        │
    ├──────────────────────────────────────────────────────>│
    │                             │    (Google Chat notify)│
    │                             │                        │
    │                             │                  4. Final Approve
```

## 📖 Example Session

```bash
# 1. First time setup
$ cool setup
🚀 Starting setup...

📧 Setting up user info...
What is your name? Jane Developer
What is your email? jane@company.com
✅ User info saved

🔗 Setting up webhooks...
GChat Review Webhook URL: https://chat.googleapis.com/v1/spaces/AAAA/messages?key=xxx
GChat Collab Webhook URL: https://chat.googleapis.com/v1/spaces/BBBB/messages?key=yyy
✅ Webhooks saved

✅ Setup complete!

# 2. Submit a review request
$ cool review request

📝 Submit Review Request to Tech Lead
=====================================

Review Title: Fix authentication bug
Description: Resolve OAuth token expiration issue
Priority:
  1. P0 - Critical
  2. P1 - High
  3. P2 - Medium
  4. P3 - Low
  5. P4 - Very Low
Select priority [3]: 1

Pull Request Links (one per line, empty line to finish):
  https://github.com/company/app/pull/456
  
Jira Ticket Links (one per line, empty line to finish):
  https://jira.company.com/SEC-789
  

📋 Preview Review Request
=========================

🔍 *New Review Request*

*Title:* Fix authentication bug
*Priority:* P0
*Submitted by:* Jane Developer (jane@company.com)
*Submitted at:* 2025-11-13 14:30:00

*Description:*
Resolve OAuth token expiration issue

*Review Links:*
• https://github.com/company/app/pull/456

*Jira Links:*
• https://jira.company.com/SEC-789

*Request ID:* `f3a2b1c4d5e6`

Do you want to submit this review request? (yes/no): yes

⏳ Submitting review request...

✅ Review request submitted successfully!

   Request ID: f3a2b1c4d5e6
   Title: Fix authentication bug
   Priority: P0
   Submitted at: 2025-11-13 14:30:00

💡 Your request has been sent to tech lead for review.
   Once approved, you can forward it to head architect using:
   cool review submit-collab f3a2b1c4d5e6

# 3. View your reviews
$ cool review histories

ID        Title                    Priority  PRs  Jira  Submitted        Collab Status  Collab Submitted
--------  -----------------------  --------  ---  ----  ---------------  -------------  ----------------
f3a2b1c4  Fix authentication bug   P0        1    1     2025-11-13 14:30 ⏳ Pending     -

Total: 1 review(s)

💡 To view details or submit to collaboration: cool review submit-collab <id>

# 4. After tech lead approval, submit to architect
$ cool review submit-collab f3a2b1c4

📋 Review Request Details
=========================

ID: f3a2b1c4d5e6
Title: Fix authentication bug
Priority: P0
Description: Resolve OAuth token expiration issue

Submitted by: Jane Developer (jane@company.com)
Submitted at: 2025-11-13 14:30:00

Pull Requests:
  • https://github.com/company/app/pull/456

Jira Tickets:
  • https://jira.company.com/SEC-789


Submit this review to head architect? (yes/no): yes

⏳ Submitting to collaboration channel...

✅ Successfully submitted to head architect!

💡 Your review request has been forwarded to the collaboration channel.
   The head architect will review and provide approval.

# 5. Hot reload during development
$ cd /Users/jane/projects/myapp

$ cool run -- go run ./cmd/http/main.go

🔥 Starting hot reload...
📄 Loading .env file from /Users/jane/projects/myapp/.env
   Watching: /Users/jane/projects/myapp
   Working directory: /Users/jane/projects/myapp
   Command: go run ./cmd/http/main.go
   Patterns: \.go$, go\.mod$, go\.sum$

Press Ctrl+C to stop
=====================================

[14:30:05] ▶️  Running: go run ./cmd/http/main.go
Server listening on :8080...

# (edit some .go files)
[14:30:42] 📝 Changes detected, restarting...
[14:30:43] ▶️  Running: go run ./cmd/http/main.go
Server listening on :8080...
```

## 🔥 Hot Reload Usage Guide

### Basic Usage

The hot reload feature automatically detects changes in your Go files and restarts your application.

#### 1. **Simple Run from Project Root**
```bash
# Navigate to project root (where go.mod is located)
cd /path/to/your/project

# Run with default settings (watches all .go files in project)
cool run -- go run .

# Run specific package
cool run -- go run ./cmd/http

# Run specific main.go
cool run -- go run cmd/http/main.go
```

#### 2. **Auto-Load .env File**
Cool CLI automatically loads `.env` file from your project root:

```bash
# Create .env in project root
cat > .env << 'EOF'
PORT=8080
DB_HOST=localhost
DB_USER=admin
DEBUG=true
EOF

# Run - automatically loads .env
cool run -- go run ./cmd/http
```

**Output:**
```
🔥 Starting hot reload...
📄 Loading .env file from /path/to/project/.env
   Watching: /path/to/project
   Working directory: /path/to/project
   Command: go run ./cmd/http
```

#### 3. **Override Environment Variables**
Flag environment variables have priority over `.env`:

```bash
# Override PORT from .env
cool run --env="PORT=3000,ENV=production" -- go run ./cmd/http
```

**Priority Order:**
1. `--env` flag (highest)
2. `.env` file
3. System environment

#### 4. **Watch Specific Directory**
Watch only specific subdirectory for better performance:

```bash
# Only watch services directory
cool run --path=./services -- go run .

# Watch multiple patterns
cool run --regex=".*\\.go$,.*\\.yaml$,.*\\.json$" -- go run .
```

### Advanced Examples

```bash
# Web server with environment
cool run --env="PORT=8080" -- go run ./cmd/api

# Microservice with custom watch patterns
cool run --regex=".*\\.go$,.*\\.yaml$" -- go run ./cmd/service

# Run tests with hot reload
cool run -- go test -v ./...

# Build and run binary
cool run -- sh -c "go build -o bin/app && ./bin/app"

# Debug mode (shows file changes)
cool run --debug -- go run ./cmd/http

# Quiet mode (no logs)
cool run --quiet -- go run ./cmd/http
```

### How It Works

1. **File Watching**: Monitors filesystem for changes (500ms polling)
2. **Working Directory**: Command runs from where you invoke `cool run` (project root with go.mod)
3. **Watch Path**: Can be different from working directory using `--path` flag
4. **Debouncing**: Waits 1 second after last change before restarting
5. **Graceful Restart**: Kills current process, starts new one

**Default Watch Patterns:**
- `\.go$` - Go source files
- `go\.mod$` - Go module file  
- `go\.sum$` - Go dependencies

**Auto-Ignored:**
- Hidden files (`.git`, `.env`, etc.)
- Hidden directories

### Common Project Structures

#### Standard Go Project
```
myapp/
├── .env                  ← Auto-loaded
├── go.mod               ← Required
├── cmd/
│   └── http/
│       └── main.go
└── internal/

# Run from project root
cd myapp
cool run -- go run ./cmd/http
```

#### Monorepo with Multiple Services
```
monorepo/
├── .env
├── go.mod
├── services/
│   ├── auth/
│   │   └── main.go
│   └── users/
│       └── main.go

# Terminal 1 - auth service
cd monorepo
cool run --env="SERVICE=auth,PORT=8081" -- go run ./services/auth

# Terminal 2 - users service  
cd monorepo
cool run --env="SERVICE=users,PORT=8082" -- go run ./services/users
```

#### Module in Subdirectory
```
project/
├── backend/
│   ├── .env
│   ├── go.mod
│   └── cmd/
│       └── api/
│           └── main.go

# Run from backend directory
cd project/backend
cool run -- go run ./cmd/api
```

### Troubleshooting

**Error: "go: go.mod file not found"**
```bash
# Make sure you're in the project root (where go.mod exists)
pwd
ls -la go.mod

# Run from correct directory
cd /path/to/project/root
cool run -- go run ./cmd/http
```

**Error: "no such file or directory"**
```bash
# Use package path, not file path
❌ cool run -- go run cmd/http/main.go  # if main.go doesn't exist at exact path
✅ cool run -- go run ./cmd/http         # Go finds main.go automatically

# Or use absolute package path
✅ cool run -- go run github.com/yourorg/project/cmd/http
```

**Changes not detected**
```bash
# Enable debug mode to see what's being watched
cool run --debug -- go run .

# Check watch patterns
cool run --regex=".*\\.go$" -- go run .
```

## 🔧 Configuration

### Editor Setup

Cool CLI supports multiple text editors for multiline input:

```bash
$ cool setup editor

📝 Choose your preferred text editor:
  1. vim
  2. nano
  3. vi
  4. emacs

Select editor [2]: 1
✅ Editor preference saved: vim
```

**Auto-Detection Priority:**
1. Configured editor (from `cool setup editor`)
2. `$EDITOR` environment variable
3. `$VISUAL` environment variable
4. First available: vim → nano → vi → emacs → code

**Tip**: The editor opens when entering review descriptions, allowing proper formatting with newlines, lists, and code blocks.

### Project Root Setup

Configure the project root for your workspace:

```bash
$ cool setup project-root

📂 Configure project root directory

Current detected paths:
  1. /Users/jane/go/src/bfi-finance
  2. /Users/jane/go/src
  3. /Users/jane/projects

Select project root or enter custom path [1]: 1
✅ Project root saved: /Users/jane/go/src/bfi-finance
```

**Usage**: When using `cool run` without `--path`, it defaults to the configured project root.

### Google Chat Webhooks

To set up Google Chat integration:

1. Create or use existing Google Chat spaces
2. Configure incoming webhooks for each space
3. Run `cool setup webhook` and paste the webhook URLs

**Required Webhooks**:
- **Review Webhook** - For tech lead notifications
- **Collab Webhook** - For head architect notifications

### Configuration Preview

View all your settings at once:

```bash
$ cool config preview

📋 Cool CLI Configuration
========================

👤 User Information:
   Name:  Jane Developer
   Email: jane@company.com

🔗 Webhooks:
   Review Webhook:  https://chat.googleapis.com/v1/spaces/AAAA/...
   Collab Webhook:  https://chat.googleapis.com/v1/spaces/BBBB/...

📝 Editor:
   Configured: vim
   Detected:   vim, nano, vi, emacs

📂 Project Root:
   /Users/jane/go/src/bfi-finance

📁 Config Location:
   ~/.cool-cli/config.json
```

## 🔥 Hot Reload Guide

### Basic Usage

```bash
# Run with default settings (watches .go, go.mod, go.sum)
cool run go run .

# Run specific package
cool run go run ./cmd/server

# Run tests with hot reload
cool run go test -v ./...

# Build and run
cool run sh -c "go build -o bin/app && ./bin/app"
```

### Advanced Usage

```bash
# Watch specific directory
cool run --path=./services go run ./services

# Custom file patterns (regex)
cool run --regex=".*\\.go$,.*\\.html$,.*\\.js$" go run .

# With environment variables
cool run --env="PORT=8080,DEBUG=true,LOG_LEVEL=debug" go run .

# Debug mode (shows detected file changes)
cool run --debug go run .

# Quiet mode (no logging)
cool run --quiet go run .

# Combine options
cool run --path=./api --env="PORT=3000" --debug go run ./api
```

### How It Works

1. **Initial Start**: Runs your command immediately
2. **File Watching**: Polls filesystem every 500ms for changes
3. **Pattern Matching**: Compares files against regex patterns
4. **Debouncing**: Waits 1 second after changes before restarting
5. **Graceful Restart**: Kills current process, starts new one

**Default Patterns**:
- `\.go$` - Go source files
- `go\.mod$` - Go module file
- `go\.sum$` - Go checksum file

**Ignored**:
- Hidden files (starting with `.`)
- Hidden directories (containing `/.` or `\.`)

### Use Cases

**Web Server Development**:
```bash
cool run --env="PORT=8080" go run ./cmd/api
```

**Microservice with Config**:
```bash
cool run --regex=".*\\.go$,.*\\.yaml$" go run ./service
```

**Integration Tests**:
```bash
cool run --path=./integration go test -v ./integration/...
```

**Multiple Services**:
```bash
# Terminal 1
cool run --env="SERVICE=auth,PORT=8081" go run ./services/auth

# Terminal 2
cool run --env="SERVICE=users,PORT=8082" go run ./services/users
```

## 💾 Data Storage

Configuration and history are stored at: `~/.cool-cli/`

### config.json
```json
{
  "gchat_review_webhook_url": "https://...",
  "gchat_collab_webhook_url": "https://...",
  "user_name": "Your Name",
  "user_email": "your@email.com",
  "preferred_editor": "vim",
  "project_root": "/Users/jane/projects/app"
}
```

### review_histories.json
```json
[
  {
    "id": "abc123...",
    "title": "Feature X",
    "priority": "high",
    "submitted_at": "2025-11-13T10:30:00Z",
    "submitted_to_collab": false,
    ...
  }
]
```

## 📚 Documentation

Comprehensive documentation available in the `docs/` directory:

- **[USER_GUIDE.md](docs/USER_GUIDE.md)** - Complete user guide
- **[REVIEW_HISTORY_SPEC.md](docs/REVIEW_HISTORY_SPEC.md)** - Complete specification (600+ lines)
- **[REVIEW_HISTORY_TESTING.md](docs/REVIEW_HISTORY_TESTING.md)** - Testing guide and scenarios
- **[QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)** - Quick command reference
- **[IMPLEMENTATION_SUMMARY.md](docs/IMPLEMENTATION_SUMMARY.md)** - Implementation details
- **[ENTITY_REFACTORING_SUMMARY.md](docs/ENTITY_REFACTORING_SUMMARY.md)** - Architecture notes
- **[REPOSITORY_INTERFACE_GUIDE.md](docs/REPOSITORY_INTERFACE_GUIDE.md)** - Repository pattern guide
- **[AI_PROMPTING_GUIDE.md](docs/AI_PROMPTING_GUIDE.md)** - AI assistance guide

## 🏗️ Clean Architecture

This project follows Clean Architecture principles with clear separation of concerns:

### Layer Structure

```
cmd/                        # Presentation Layer
├── Commands handle user input and output
├── Thin layer, delegates to use cases
└── CLI-specific concerns only

internal/domain/           # Business Logic Layer
├── entity/               # Core domain models
│   ├── review_history.go
│   └── hot_reload.go
├── repository/           # Port interfaces
│   ├── review_history.go
│   └── hot_reload.go
└── usecase/              # Business use cases
    ├── review.go
    ├── g_chat.go
    └── hot_reload.go

internal/infrastructure/   # Infrastructure Layer
└── repository/           # Adapter implementations
    ├── review_history.go  (JSON storage)
    ├── file_watcher.go    (Polling-based)
    └── process_runner.go  (Command execution)

internal/pkg/             # Shared Utilities
├── common/               # Common helpers
├── logger/               # Logging
└── table/                # Table display

config/                   # Configuration
└── config.go            # Config management
```

### Design Principles

1. **Dependency Rule**: Dependencies point inward
   - Presentation → Use Cases → Entities
   - Infrastructure → Use Cases (through interfaces)

2. **Interface Segregation**: Small, focused interfaces
   - `FileWatcher`: Watch, Stop
   - `ProcessRunner`: Start, Stop, IsRunning, Info

3. **Dependency Inversion**: Depend on abstractions
   - Domain defines interfaces (ports)
   - Infrastructure provides implementations (adapters)

4. **Single Responsibility**: Each layer has one reason to change
   - Domain: Business logic changes
   - Infrastructure: Technology changes
   - Presentation: UI/UX changes

### Benefits

- ✅ **Testability**: Easy to mock interfaces
- ✅ **Maintainability**: Clear boundaries between layers
- ✅ **Flexibility**: Swap implementations without touching business logic
- ✅ **Scalability**: Add new features without breaking existing code

## 🔧 Priority Levels

Cool CLI uses a **P0-P4 priority system** for review requests:

- **P0 - Critical**: Highest priority, urgent issues requiring immediate attention
  - Production down, security vulnerabilities, critical bugs
  - Expected review time: Within hours
  
- **P1 - High**: High priority, important features or significant bugs
  - Major features, important bug fixes, performance issues
  - Expected review time: Within 1 day
  
- **P2 - Medium** (Default): Normal priority for standard changes
  - Regular features, minor bug fixes, refactoring
  - Expected review time: Within 2-3 days
  
- **P3 - Low**: Low priority, nice-to-have changes
  - Code cleanup, documentation updates, minor improvements
  - Expected review time: Within 1 week
  
- **P4 - Very Low**: Lowest priority, optional changes
  - Cosmetic changes, experimental features, future considerations
  - Expected review time: When available

**Selection**: When submitting a review request, you'll see a menu:
```
Priority:
  1. P0 - Critical
  2. P1 - High
  3. P2 - Medium
  4. P3 - Low
  5. P4 - Very Low
Select priority [3]:
```
Simply type the number (1-5) corresponding to your priority level.

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
# Clone repository
git clone https://github.com/yatbfi/cool.git
cd cool

# Install dependencies
go mod tidy

# Build
go build -o bin/cool .

# Run
./bin/cool --help
```

### Hot Reload Development

Use cool itself for development:

```bash
# Build and run with hot reload
cool run go run .

# Run tests with hot reload
cool run go test -v ./...

# Custom patterns for templates/configs
cool run --regex=".*\\.go$,.*\\.yaml$,.*\\.json$" go run .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Project Structure

```
cool/
├── cmd/                     # Command implementations (Presentation)
│   ├── root.go             # Root command
│   ├── setup*.go           # Setup commands
│   ├── review*.go          # Review commands
│   ├── run.go              # Hot reload command
│   └── ...
├── config/                  # Configuration management
│   └── config.go
├── internal/
│   ├── domain/             # Domain layer (Business logic)
│   │   ├── entity/        # Domain models
│   │   ├── repository/    # Repository interfaces (ports)
│   │   └── usecase/       # Business use cases
│   ├── infrastructure/     # Infrastructure layer (Adapters)
│   │   └── repository/    # Repository implementations
│   └── pkg/               # Internal packages
│       ├── common/        # Common utilities
│       ├── logger/        # Logging
│       └── table/         # Table rendering
├── docs/                   # Documentation
├── bin/                    # Compiled binaries
├── main.go                 # Entry point
├── go.mod                  # Go modules
├── go.sum                  # Go dependencies checksum
├── Makefile                # Build automation
└── README.md               # This file
```

### Adding New Features

Follow clean architecture when adding features:

1. **Define Domain Entity** (if needed)
   ```go
   // internal/domain/entity/new_feature.go
   package entity
   
   type NewFeature struct {
       ID   string
       Name string
   }
   ```

2. **Define Repository Interface**
   ```go
   // internal/domain/repository/new_feature.go
   package repository
   
   type NewFeatureRepository interface {
       Save(feature *entity.NewFeature) error
       FindByID(id string) (*entity.NewFeature, error)
   }
   ```

3. **Create Use Case**
   ```go
   // internal/domain/usecase/new_feature.go
   package usecase
   
   type NewFeatureUsecase struct {
       repo repository.NewFeatureRepository
   }
   
   func (uc *NewFeatureUsecase) DoSomething() error {
       // Business logic here
   }
   ```

4. **Implement Repository**
   ```go
   // internal/infrastructure/repository/new_feature.go
   package repository
   
   type newFeatureRepository struct {
       // implementation details
   }
   
   func (r *newFeatureRepository) Save(feature *entity.NewFeature) error {
       // save logic
   }
   ```

5. **Create Command**
   ```go
   // cmd/new_feature.go
   package cmd
   
   func NewNewFeatureCmd() *cobra.Command {
       return &cobra.Command{
           Use: "new-feature",
           RunE: func(cmd *cobra.Command, args []string) error {
               // Create dependencies
               repo := repository.NewNewFeatureRepository()
               uc := usecase.NewNewFeatureUsecase(repo)
               
               // Execute use case
               return uc.DoSomething()
           },
       }
   }
   ```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. **Follow Clean Architecture**: Keep layers separated
   - Domain layer must not import infrastructure
   - Use dependency injection
   - Define interfaces in domain layer

2. **Write Tests**: Add tests for new features
   ```bash
   go test ./internal/domain/usecase/...
   ```

3. **Update Documentation**: 
   - Update README.md for user-facing changes
   - Add docs in `docs/` for complex features
   - Add code comments for public APIs

4. **Follow Go Best Practices**:
   - Use `gofmt` for formatting
   - Run `go vet` for static analysis
   - Follow [Effective Go](https://go.dev/doc/effective_go)

5. **Commit Messages**:
   ```
   feat: add hot reload feature
   fix: resolve editor detection on Windows
   docs: update README with hot reload guide
   refactor: extract file watcher to infrastructure layer
   ```

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following guidelines
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📝 License

[Add your license here]

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [promptui](https://github.com/manifoldco/promptui) for interactive prompts
- Inspired by [Air](https://github.com/cosmtrek/air) and [Watch](https://github.com/tanyudii/watch) for hot reload
- Thanks to all contributors

## 📞 Support

For issues, questions, or suggestions:
- 🐛 [Create an issue](https://github.com/yatbfi/cool/issues) on GitHub
- 📖 Check the documentation in `docs/`
- 📋 Review the [QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)

## 🚦 Status

- ✅ **Production Ready**
- ✅ **Clean Architecture**
- ✅ **Comprehensive Documentation**
- ✅ **Hot Reload Support**
- ✅ **Thread-Safe**
- ✅ **Well-Tested**

## 📈 Roadmap

### Current Version (3.0.0)
- ✅ Review request submission with editor
- ✅ Preview and edit before submit
- ✅ P0-P4 priority system
- ✅ History tracking with filters
- ✅ Collaboration forwarding
- ✅ Google Chat integration
- ✅ Beautiful table display
- ✅ Thread-safe JSON storage
- ✅ Clean architecture implementation
- ✅ Hot reload functionality
- ✅ Editor auto-detection and setup
- ✅ Project root configuration
- ✅ Configuration preview

### Future Enhancements
- [ ] Edit/delete review entries
- [ ] Search and advanced filter functionality
- [ ] Manual approval status tracking
- [ ] Comments and notes on reviews
- [ ] Export to CSV/Excel
- [ ] GitHub API integration (auto-fetch PR details)
- [ ] Jira API integration (auto-fetch ticket details)
- [ ] Statistics dashboard with metrics
- [ ] Email notifications
- [ ] Slack integration
- [ ] Web UI interface
- [ ] Hot reload for non-Go projects
- [ ] Watch exclusion patterns
- [ ] Custom restart commands
- [ ] Process health checks

---

**Made with ❤️ for developers who value efficient workflows**

**Version**: 3.0.0  
**Last Updated**: January 2025  
**Status**: Production Ready ✅

