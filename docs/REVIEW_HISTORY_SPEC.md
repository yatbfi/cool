# Review History System - Complete Specification

## Overview

The Review History System adalah fitur untuk mengelola workflow review request dari developer ke tech lead dan head architect. System ini menyimpan history semua review request dan memungkinkan tracking status submission ke collaboration channel.

## Use Case Flow

```
┌──────────────┐
│  Developer   │
└──────┬───────┘
       │
       │ 1. Submit Review Request
       ▼
┌──────────────────────┐
│    Tech Lead         │
│  (Review Webhook)    │
└──────┬───────────────┘
       │
       │ 2. Approve
       │ 3. Forward to Collaboration
       ▼
┌──────────────────────┐
│  Head Architect      │
│ (Collab Webhook)     │
└──────────────────────┘
```

## Components

### 1. Domain Layer

#### Repository Interface: `ReviewHistoryRepository`

**File**: `internal/domain/repository/review_history.go`

```go
type ReviewHistoryEntry struct {
    ID                   string
    Title                string
    Description          string
    Priority             string
    ReviewLinks          []string
    JiraLinks            []string
    SubmittedBy          string
    SubmittedByEmail     string
    SubmittedAt          time.Time
    SubmittedToCollab    bool
    SubmittedToCollabAt  *time.Time
    SubmittedToCollabBy  string
    ApprovedByTechLead   bool
    ApprovedByArchitect  bool
    Notes                string
}
```

**Methods**:
- `Save(ctx, entry)` - Save new review history
- `Update(ctx, entry)` - Update existing history
- `FindByID(ctx, id)` - Get history by ID
- `FindAll(ctx)` - Get all histories
- `FindByCollabStatus(ctx, submitted)` - Filter by collab status
- `Delete(ctx, id)` - Delete history

#### Usecase: `Review`

**File**: `internal/domain/usecase/review.go`

**Methods**:
- `SubmitReviewRequest(ctx, req)` - Submit new review to tech lead
- `GetHistories(ctx, filter)` - Get histories with filter
- `GetHistoryByID(ctx, id)` - Get specific history
- `SubmitToCollaboration(ctx, historyID)` - Forward to head architect
- `SendToGChat(ctx, webhookURL, message)` - Send message to GChat

**Filter Types**:
- `HistoryFilterAll` - All histories
- `HistoryFilterPending` - Not submitted to collab
- `HistoryFilterCompleted` - Already submitted to collab

### 2. Infrastructure Layer

#### Repository Implementation

**File**: `internal/infrastructure/repository/review_history.go`

**Storage**: JSON file at `~/.cool-cli/review_histories.json`

**Features**:
- Thread-safe read/write with mutex
- Automatic file initialization
- JSON marshaling/unmarshaling

### 3. Command Layer

#### Commands

##### `cool review request`

Submit new review request to tech lead.

**Flow**:
1. Prompt user for input (title, description, priority, links)
2. Validate webhook configuration
3. Generate unique ID
4. Save to history
5. Send to GChat review webhook
6. Display success message with ID

**Input Fields**:
- Title (required)
- Description (optional)
- Priority: low/medium/high/critical (default: medium)
- Pull Request Links (multiple)
- Jira Ticket Links (multiple)

**Output**:
```
✅ Review request submitted successfully!

   Request ID: abc123...
   Title: Feature X Review
   Priority: high
   Submitted at: 2025-01-15 10:30:00

💡 Your request has been sent to tech lead for review.
   Once approved, you can forward it to head architect using:
   cool review submit-collab abc123
```

##### `cool review histories`

Display review history table.

**Usage**:
```bash
cool review histories              # Show all
cool review histories --pending    # Show pending only
cool review histories --completed  # Show completed only
```

**Table Columns**:
- ID (truncated to 8 chars)
- Title (truncated to 40 chars)
- Priority
- PRs (count)
- Jira (count)
- Submitted (timestamp)
- Collab Status (✅/⏳)
- Collab Submitted (timestamp if submitted)

**Output Example**:
```
ID       Title                     Priority  PRs  Jira  Submitted        Collab Status  Collab Submitted
-------- ------------------------- --------  ---  ----  ---------------  -------------  ----------------
abc12345 Feature X Review          high      2    1     2025-01-15 10:30 ⏳ Pending     -
def67890 Bug Fix Review            medium    1    2     2025-01-14 15:20 ✅ Submitted   2025-01-14 16:00

Total: 2 review(s)
```

##### `cool review submit-collab [id]`

Forward approved review to head architect.

**Usage**:
```bash
cool review submit-collab abc123           # Submit specific review
cool review submit-collab --list           # Show all reviews
cool review submit-collab --list --pending # Show pending only
```

**Flow**:
1. Validate collab webhook configuration
2. Get history by ID
3. Check if already submitted
4. Display review details
5. Ask for confirmation
6. Send to GChat collab webhook
7. Update history entry
8. Display success message

**Confirmation Prompt**:
```
📋 Review Request Details
=========================

ID: abc123...
Title: Feature X Review
Priority: high
Description: Add new feature X

Submitted by: John Doe (john@example.com)
Submitted at: 2025-01-15 10:30:00

Pull Requests:
  • https://github.com/org/repo/pull/123

Jira Tickets:
  • https://jira.example.com/PROJ-456

Submit this review to head architect? (yes/no):
```

### 4. Configuration

#### Config Structure

**File**: `config/config.go`

```go
type Config struct {
    GChatReviewWebhookURL string `json:"gchat_review_webhook_url"`
    GChatCollabWebhookURL string `json:"gchat_collab_webhook_url"`
    UserName              string `json:"user_name"`
    UserEmail             string `json:"user_email"`
}
```

**Storage**: `~/.cool-cli/config.json`

#### Setup Commands

##### `cool setup email`

Configure user name and email.

##### `cool setup webhook`

Configure both webhook URLs:
- Review webhook (for tech lead)
- Collaboration webhook (for head architect)

### 5. Validation

#### Pre-run Validation

**File**: `cmd/base.go`

**Checks**:
1. Environment validation (OS, Go binary)
2. Skip for setup/update commands
3. User setup validation (name, email)
4. Command-specific validation:
   - `request` command: Requires review webhook
   - `submit-collab` command: Requires collab webhook

**Error Messages**:
```
⚠️  GChat review webhook URL is not configured.
Please run the setup command to configure it:
   cool setup webhook
```

## Google Chat Integration

### Message Format - Review Request

```
🔍 *New Review Request*

*Title:* Feature X Review
*Priority:* high
*Submitted by:* John Doe (john@example.com)
*Submitted at:* 2025-01-15 10:30:00

*Description:*
Add new feature X with proper tests

*Review Links:*
• https://github.com/org/repo/pull/123
• https://github.com/org/repo/pull/124

*Jira Links:*
• https://jira.example.com/PROJ-456

*Request ID:* `abc123...`
```

### Message Format - Collaboration

```
🚀 *Review Request for Head Architect Approval*

*Title:* Feature X Review
*Priority:* high
*Originally submitted by:* John Doe (john@example.com)
*Tech Lead Approved:* ✅
*Forwarded at:* 2025-01-15 14:00:00

*Description:*
Add new feature X with proper tests

*Review Links:*
• https://github.com/org/repo/pull/123

*Jira Links:*
• https://jira.example.com/PROJ-456

*Request ID:* `abc123...`

_Please review and approve._
```

## Data Storage

### File Structure

```
~/.cool-cli/
├── config.json              # User config (webhooks, name, email)
└── review_histories.json    # All review histories
```

### Review History JSON Format

```json
[
  {
    "id": "abc123...",
    "title": "Feature X Review",
    "description": "Add new feature X",
    "priority": "high",
    "review_links": [
      "https://github.com/org/repo/pull/123"
    ],
    "jira_links": [
      "https://jira.example.com/PROJ-456"
    ],
    "submitted_by": "John Doe",
    "submitted_by_email": "john@example.com",
    "submitted_at": "2025-01-15T10:30:00Z",
    "submitted_to_collab": false,
    "submitted_to_collab_at": null,
    "submitted_to_collab_by": "",
    "approved_by_tech_lead": false,
    "approved_by_architect": false,
    "notes": ""
  }
]
```

## Error Handling

### Common Errors

1. **Webhook Not Configured**
   - Show clear message with setup instructions
   - Return error before attempting operation

2. **Review Not Found**
   - Check ID exists in history
   - Return user-friendly error

3. **Already Submitted**
   - Prevent duplicate submission to collab
   - Show when it was submitted

4. **Network Errors**
   - Handle GChat webhook failures gracefully
   - Return descriptive error messages

## Future Enhancements

### Possible Features

1. **Approval Tracking**
   - Mark reviews as approved by tech lead
   - Mark reviews as approved by architect
   - Add approval timestamps

2. **Comments/Notes**
   - Add notes to history entries
   - Track feedback from reviewers

3. **Search & Filter**
   - Search by title/description
   - Filter by priority
   - Filter by date range
   - Filter by submitter

4. **Export**
   - Export history to CSV
   - Generate reports

5. **Reminders**
   - Reminder for pending reviews
   - Notification when reviews are stuck

6. **Integration**
   - GitHub API integration
   - Jira API integration
   - Auto-fetch PR details

7. **Statistics**
   - Review time metrics
   - Approval rates
   - Most active reviewers

## Testing Recommendations

### Unit Tests

1. Repository tests
   - CRUD operations
   - Concurrent access
   - File corruption handling

2. Usecase tests
   - Review submission flow
   - History filtering
   - Collaboration submission
   - Error scenarios

3. Command tests
   - Input validation
   - User interaction
   - Output formatting

### Integration Tests

1. End-to-end flow
   - Submit request → View history → Submit to collab
   - Webhook mock server testing
   - File system operations

### Manual Testing

1. Happy path testing
2. Error scenario testing
3. Edge cases (empty lists, long strings, special characters)
4. Concurrent usage testing

## Usage Examples

### Complete Workflow Example

```bash
# 1. Setup (first time)
cool setup email
cool setup webhook

# 2. Submit review request
cool review request
# Enter details interactively...

# 3. View all histories
cool review histories

# 4. View pending only
cool review histories --pending

# 5. Submit to collaboration (after tech lead approval)
cool review submit-collab abc123

# 6. List with submit-collab
cool review submit-collab --list --pending

# 7. View completed submissions
cool review histories --completed
```

## Architecture Benefits

1. **Clean Architecture**
   - Clear separation of concerns
   - Easy to test and maintain
   - Flexible for future changes

2. **Repository Pattern**
   - Abstract data storage
   - Easy to switch storage (JSON → DB)
   - Mockable for testing

3. **Usecase Layer**
   - Business logic isolation
   - Reusable across commands
   - Clear dependencies

4. **Command Pattern**
   - Consistent UX
   - Shared validation
   - Easy to add new commands

## Conclusion

This specification provides a complete review history system that:
- Tracks all review requests
- Integrates with Google Chat
- Provides clear workflow for developers
- Maintains audit trail
- Supports future enhancements

The system is designed to be simple, reliable, and developer-friendly.

