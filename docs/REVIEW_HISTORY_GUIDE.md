# Review History Feature

## Overview

The Review History feature allows developers to track their code review requests and forward approved reviews to the collaboration channel for head architect approval.

## Workflow

```
┌─────────────┐       ┌─────────────┐       ┌──────────────────┐
│  Developer  │──────▶│  Tech Lead  │──────▶│  Head Architect  │
└─────────────┘       └─────────────┘       └──────────────────┘
     │                      │                        │
     │ review request       │ approve               │
     ▼                      ▼                        ▼
[Review Channel]      [Decision]         [Collaboration Channel]
     │                      │                        │
     └──────────────────────┴────────────────────────┘
              History tracked in ~/.cool-cli/
```

### Step by Step

1. **Developer submits review request**
   ```bash
   cool review request
   ```
   - Sends review to Tech Lead channel
   - Saves history locally

2. **Tech Lead reviews** (external step)
   - Tech Lead reviews the PR
   - Provides feedback or approves

3. **Developer forwards approved review**
   ```bash
   cool review submit-collab
   ```
   - Shows pending reviews
   - Forwards to Head Architect channel
   - Updates history status

## Commands

### `cool review request`

Submit a new review request to the Tech Lead channel.

**Usage:**
```bash
cool review request              # Interactive mode
cool review request --dry-run    # Preview without sending
```

**Interactive prompts:**
- Title
- PR links (Service + URL)
- Jira links
- Description
- Priority (P0, P1, P2, P3)

**After submission:**
- Message sent to GChat review channel
- Entry saved to `~/.cool-cli/review_histories.json`
- Unique ID generated for tracking

**Example:**
```bash
$ cool review request
? Changes title: Add user authentication
? Add PR link? (y/N) y
? PR Service (e.g., LSS, LGS, LTS, LTW, LPW): LSS
? PR Link: https://github.com/org/repo/pull/123
? Add PR link? (y/N) n
? Add Jira link? (y/N) y
? Jira Link: https://jira.company.com/browse/PROJ-123
? Add Jira link? (y/N) n
? Description: Implement JWT-based authentication
? Priority: P1
✅ Review request sent successfully! (ID: a1b2c3d4)
```

---

### `cool review histories`

Display all review request history with submission status.

**Usage:**
```bash
cool review histories              # Show all history
cool review histories --pending    # Show only pending (not submitted to collab)
cool review histories --completed  # Show only submitted to collab
```

**Output:**

Table showing:
- **ID**: Short ID (first 8 chars)
- **Title**: Review title
- **Priority**: P0, P1, P2, or P3
- **PRs**: Number of pull requests
- **Jira**: Number of Jira tickets
- **Submitted**: When review was submitted to Tech Lead
- **Collab Status**: ⏳ Pending or ✅ Submitted
- **Collab Submitted**: When submitted to collaboration (if applicable)

**Example:**
```bash
$ cool review histories

ID        Title                      Priority  PRs  Jira  Submitted            Collab Status  Collab Submitted
--------  -----                      --------  ---  ----  ---------            -------------  ----------------
a1b2c3d4  [P1] Add user auth         P1        2    1     2025-11-13 10:30     ✅ Submitted   2025-11-13 14:20
e5f6g7h8  [P0] Fix payment gateway   P0        1    2     2025-11-13 09:15     ⏳ Pending     -
i9j0k1l2  [P3] Update API docs       P3        1    0     2025-11-12 16:45     ⏳ Pending     -

Total: 3 review(s)
```

**Filter examples:**
```bash
# Show only pending reviews
$ cool review histories --pending

# Show only completed reviews
$ cool review histories --completed
```

---

### `cool review submit-collab`

Forward an approved review to the collaboration channel for head architect approval.

**Usage:**
```bash
cool review submit-collab                    # Interactive selection
cool review submit-collab --id a1b2c3d4      # Submit specific review by ID
cool review submit-collab --latest           # Submit latest pending review
```

**Interactive Flow:**

1. Shows table of pending reviews
2. Prompts for selection
3. Shows preview of selected review
4. Asks for confirmation
5. Submits to collaboration channel
6. Updates history status

**Example (Interactive):**
```bash
$ cool review submit-collab

Select review to submit to collaboration channel:

No.  ID        Title                      Priority  PRs  Jira  Submitted
---  --------  -----                      --------  ---  ----  ---------
1    e5f6g7h8  [P0] Fix payment gateway   P0        1    2     2025-11-13 09:15
2    i9j0k1l2  [P3] Update API docs       P3        1    0     2025-11-12 16:45

? Select review [1-2]: 1

============================================================
📋 Review Preview - Will be sent to Collaboration Channel
============================================================

📝 Title: [P0] Fix payment gateway
📄 Description: Critical bug fix for payment processing
🏷️  Priority: P0
📅 Originally Submitted: 2025-11-13 09:15
👤 Submitted By: John Doe (john.doe@company.com)

🔗 Pull Requests:
  1. LSS: https://github.com/org/repo/pull/456

🎫 Jira Links:
  1. https://jira.company.com/browse/PROJ-456
  2. https://jira.company.com/browse/PROJ-457

============================================================
? Submit this review to collaboration channel (y/N) y

⏳ Submitting to collaboration channel...
✅ Review submitted to collaboration channel!
```

**Example (Direct by ID):**
```bash
$ cool review submit-collab --id e5f6g7h8
📋 Review Preview...
✅ Review submitted to collaboration channel!
```

**Example (Latest):**
```bash
$ cool review submit-collab --latest
📋 Submitting latest pending review: [P0] Fix payment gateway
✅ Review submitted to collaboration channel!
```

**Validation:**
- Must have GChat collaboration webhook URL configured
- Review must exist in history
- Review must not be already submitted to collaboration

---

## Data Storage

### Location
`~/.cool-cli/review_histories.json`

### Format
JSON array of review entries:

```json
[
  {
    "id": "1731488400-123456789",
    "title": "[P1] Add user authentication",
    "description": "Implement JWT-based auth",
    "review_links": [
      {
        "Service": "LSS",
        "PullRequestURL": "https://github.com/org/repo/pull/123"
      }
    ],
    "jira_links": [
      "https://jira.company.com/browse/PROJ-123"
    ],
    "priority": "P1",
    "submitted_at": "2025-11-13T10:30:00+07:00",
    "submitted_to_collab": false,
    "submitted_by": "John Doe (john.doe@company.com)"
  }
]
```

### File Permissions
- **Mode**: `0600` (owner read/write only)
- **Location**: Same directory as config

---

## Configuration

### Required Webhooks

Both webhook URLs must be configured for full functionality:

```bash
# Setup review webhook (for Tech Lead channel)
cool setup webhook

# This will prompt for both:
# 1. GChat review webhook URL (required for `cool review request`)
# 2. GChat collaboration webhook URL (required for `cool review submit-collab`)
```

### Validation

The tool validates webhook configuration before running commands:

- `cool review request` → Requires review webhook URL
- `cool review submit-collab` → Requires collaboration webhook URL

If missing, you'll see:
```
⚠️  GChat review webhook URL is not configured.
Please run the setup command to configure it:
   "cool setup webhook"

Error: missing GChat review webhook URL
```

---

## Message Format

### Review Request Message (to Tech Lead)

```
*[P1] Add user authentication*

_Implement JWT-based authentication_

*Pull Requests:*
• LSS → <https://github.com/org/repo/pull/123|#123>

*Jira Links:*
• <https://jira.company.com/browse/PROJ-123|PROJ-123>
```

### Collaboration Message (to Head Architect)

Same format but with additional context:

```
*[APPROVED BY TECH LEAD] [P1] Add user authentication*

_Implement JWT-based authentication_

*Pull Requests:*
• LSS → <https://github.com/org/repo/pull/123|#123>

*Jira Links:*
• <https://jira.company.com/browse/PROJ-123|PROJ-123>

*Originally Submitted:* 2025-11-13 10:30
*Submitted by:* John Doe (john.doe@company.com)
```

---

## Error Handling

### Common Scenarios

**No pending reviews:**
```bash
$ cool review submit-collab
✅ No pending reviews to submit to collaboration.
   All your reviews have been submitted!
```

**Review not found:**
```bash
$ cool review submit-collab --id xyz123
Error: review with ID 'xyz123' not found
```

**Already submitted:**
```bash
$ cool review submit-collab --id a1b2c3d4
Error: review already submitted to collaboration at 2025-11-13 14:20
```

**History file corrupted:**
```
Error: failed to parse history file (backup created at ~/.cool-cli/review_histories.json.backup)
```

---

## Tips & Best Practices

### 1. Use Dry-Run for Preview
Always preview your review before submitting:
```bash
cool review request --dry-run
```

### 2. Check Pending Reviews Regularly
```bash
cool review histories --pending
```

### 3. Submit Latest Quickly
For the most recent pending review:
```bash
cool review submit-collab --latest
```

### 4. Use Full ID for Scripts
When automating, use the full ID:
```bash
cool review submit-collab --id 1731488400-123456789
```

### 5. Filter by Status
Track your workflow:
```bash
# What's pending?
cool review histories --pending

# What's been completed?
cool review histories --completed
```

---

## Architecture

This feature follows **Hexagonal Architecture** (Ports & Adapters):

```
┌─────────────────────────────────────────────────┐
│                   Presentation                   │
│  (cmd/review_*.go - CLI commands)               │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│                   Domain/Use Cases               │
│  (usecase/review.go - Business logic)          │
│  - RequestReview()                               │
│  - GetHistories()                                │
│  - SubmitToCollaboration()                      │
└────────────────┬────────────────────────────────┘
                 │
      ┌──────────┴──────────┐
      │                     │
┌─────▼─────┐      ┌───────▼────────┐
│  GChat    │      │   History      │
│  Port     │      │   Repository   │
│           │      │   Port         │
└─────┬─────┘      └───────┬────────┘
      │                    │
┌─────▼─────┐      ┌───────▼────────┐
│  GChat    │      │   File-based   │
│  Adapter  │      │   Adapter      │
└───────────┘      └────────────────┘
```

**Benefits:**
- Business logic isolated from infrastructure
- Easy to test
- Easy to swap implementations (e.g., database instead of files)
- Clear separation of concerns

---

## Future Enhancements

See [REVIEW_HISTORY_SPEC.md](./REVIEW_HISTORY_SPEC.md) for detailed future ideas:

- Review status tracking (pending_tech_lead, approved, etc.)
- Comments/notes on reviews
- Advanced filtering and search
- Export to CSV/Markdown
- Statistics and analytics
- Auto-cleanup and archiving
- Git integration for auto-detection

---

## Troubleshooting

### History file not created
**Solution:** Submit your first review with `cool review request`

### Can't see my review in histories
**Solution:** Ensure the review was submitted successfully (look for "Review request sent successfully!" message)

### Wrong webhook being used
**Solution:** Check your configuration:
```bash
cat ~/.cool-cli/config.json
```

Ensure both webhooks are set:
- `gchat_review_webhook_url` - for review requests
- `gchat_collab_webhook_url` - for collaboration submissions

### Table formatting looks weird
**Solution:** Ensure your terminal width is at least 120 characters

---

## Support

For issues or questions:
1. Check this documentation
2. Run commands with `--help` flag
3. Check your configuration file
4. Review the spec document: `docs/REVIEW_HISTORY_SPEC.md`

