# Cool CLI - Developer Review Workflow Tool 🚀

A powerful command-line tool to streamline the code review workflow from developer to tech lead to head architect, with seamless Google Chat integration.

## ✨ Features

- 📝 **Review Request Submission** - Submit review requests to tech lead with all necessary details
- 👀 **Preview Before Submit** - Preview formatted message before sending to avoid mistakes
- 🎯 **P0-P4 Priority System** - Clear priority levels from Critical to Very Low
- 📊 **History Tracking** - Track all your review requests with comprehensive history
- 🔄 **Collaboration Forwarding** - Forward approved reviews to head architect
- 💬 **Google Chat Integration** - Automatic notifications to review channels
- 🎨 **Beautiful Table Display** - Clean and organized view of your reviews
- ⚙️ **Easy Configuration** - Simple setup wizard for user info and webhooks
- 🔁 **Input Retry Logic** - Forgiving input validation with retry on errors
- 🛡️ **Thread-Safe** - Safe concurrent operations with mutex protection

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
# Run setup wizard
cool setup

# Or configure separately
cool setup email     # Configure your name and email
cool setup webhook   # Configure Google Chat webhooks
```

### Basic Usage

```bash
# Submit a review request
cool review request

# View all review history
cool review histories

# View pending reviews only
cool review histories --pending

# Submit to collaboration (after tech lead approval)
cool review submit-collab <review-id>

# List reviews with submit-collab
cool review submit-collab --list --pending
```

## 📋 Commands

### Setup Commands

| Command | Description |
|---------|-------------|
| `cool setup` | Run full setup wizard |
| `cool setup email` | Configure user name and email |
| `cool setup webhook` | Configure Google Chat webhook URLs |

### Review Commands

| Command | Description |
|---------|-------------|
| `cool review request` | Submit new review request to tech lead |
| `cool review histories` | Display all review history |
| `cool review histories --pending` | Show pending reviews only |
| `cool review histories --completed` | Show completed reviews only |
| `cool review submit-collab <id>` | Submit review to head architect |
| `cool review submit-collab --list` | List all reviews |
| `cool review submit-collab --list --pending` | List pending reviews |

### Other Commands

| Command | Description |
|---------|-------------|
| `cool update` | Update Cool CLI to latest version |
| `cool completion` | Generate shell completion scripts |

## 🎯 Workflow

```
Developer                     Tech Lead              Head Architect
    │                             │                        │
    │ 1. cool review request      │                        │
    ├────────────────────────────>│                        │
    │    (Google Chat notify)     │                        │
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
```

## 🏗️ Architecture

```
internal/
├── domain/
│   ├── entity/          # Domain models (ReviewHistoryEntry)
│   ├── repository/      # Repository interfaces
│   └── usecase/        # Business logic (Review, GChat)
├── infrastructure/
│   └── repository/      # Repository implementations (JSON storage)
└── pkg/
    ├── common/          # Common utilities
    ├── logger/          # Logging utilities
    └── table/           # Custom table rendering
```

**Design Patterns**:
- Clean Architecture
- Repository Pattern
- Dependency Injection
- Command Pattern (Cobra)

## 💾 Data Storage

Configuration and history are stored at: `~/.cool-cli/`

### config.json
```json
{
  "gchat_review_webhook_url": "https://...",
  "gchat_collab_webhook_url": "https://...",
  "user_name": "Your Name",
  "user_email": "your@email.com"
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

- **[REVIEW_HISTORY_SPEC.md](docs/REVIEW_HISTORY_SPEC.md)** - Complete specification (600+ lines)
- **[REVIEW_HISTORY_TESTING.md](docs/REVIEW_HISTORY_TESTING.md)** - Testing guide and scenarios
- **[QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)** - Quick command reference
- **[IMPLEMENTATION_SUMMARY.md](docs/IMPLEMENTATION_SUMMARY.md)** - Implementation details
- **[ENTITY_REFACTORING_SUMMARY.md](docs/ENTITY_REFACTORING_SUMMARY.md)** - Architecture notes

## 🔧 Configuration

### Google Chat Webhooks

To set up Google Chat integration:

1. Create or use existing Google Chat spaces
2. Configure incoming webhooks for each space
3. Run `cool setup webhook` and paste the webhook URLs

**Required Webhooks**:
- **Review Webhook** - For tech lead notifications
- **Collab Webhook** - For head architect notifications

### Priority Levels

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

### Running Tests

```bash
go test ./...
```

### Project Structure

```
cool/
├── cmd/                 # Command implementations
├── config/              # Configuration management
├── internal/
│   ├── domain/          # Domain layer
│   ├── infrastructure/  # Infrastructure layer
│   └── pkg/             # Internal packages
├── docs/                # Documentation
├── bin/                 # Compiled binaries
├── main.go              # Entry point
├── go.mod               # Go modules
└── README.md            # This file
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. Follow clean architecture principles
2. Write tests for new features
3. Update documentation
4. Follow Go best practices
5. Add helpful comments

## 📝 License

[Add your license here]

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Inspired by developer workflow optimization needs
- Thanks to all contributors

## 📞 Support

For issues, questions, or suggestions:
- Create an issue on GitHub
- Check the documentation in `docs/`
- Review the quick reference guide

## 🚦 Status

- ✅ **Production Ready**
- ✅ **Clean Architecture**
- ✅ **Comprehensive Documentation**
- ✅ **Thread-Safe**
- ✅ **Well-Tested**

## 📈 Roadmap

### Current Version (2.1.0)
- ✅ Review request submission
- ✅ Preview before submit
- ✅ P0-P4 priority system with retry validation
- ✅ History tracking with filters
- ✅ Collaboration forwarding
- ✅ Google Chat integration
- ✅ Beautiful table display with custom package
- ✅ Thread-safe JSON storage
- ✅ Clean architecture implementation

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

---

**Made with ❤️ for developers who value efficient workflows**

**Version**: 2.1.0  
**Last Updated**: November 13, 2025  
**Status**: Production Ready ✅

