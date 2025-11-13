# Cool CLI - Quick Reference Guide

Quick reference untuk menggunakan Cool CLI review system.

## 🚀 Quick Start

```bash
# 1. Setup (first time only)
cool setup

# 2. Submit review request
cool review request

# 3. View your reviews
cool review histories

# 4. Submit to collaboration (after tech lead approval)
cool review submit-collab <review-id>
```

## 📋 Command Reference

### Setup Commands

| Command | Description | Example |
|---------|-------------|---------|
| `cool setup` | Run full setup | `cool setup` |
| `cool setup email` | Configure name & email | `cool setup email` |
| `cool setup webhook` | Configure webhook URLs | `cool setup webhook` |

### Review Commands

| Command | Description | Example |
|---------|-------------|---------|
| `cool review request` | Submit new review | `cool review request` |
| `cool review histories` | View all reviews | `cool review histories` |
| `cool review histories --pending` | View pending reviews | `cool review histories --pending` |
| `cool review histories --completed` | View completed reviews | `cool review histories --completed` |
| `cool review submit-collab <id>` | Submit to head architect | `cool review submit-collab abc123` |
| `cool review submit-collab --list` | List all reviews | `cool review submit-collab --list` |
| `cool review submit-collab --list --pending` | List pending reviews | `cool review submit-collab --list --pending` |

## 🔧 Configuration

Config files stored at: `~/.cool-cli/`

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
    "id": "unique-id",
    "title": "Review Title",
    "priority": "high",
    "submitted_at": "2025-01-15T10:30:00Z",
    "submitted_to_collab": false,
    ...
  }
]
```

## 📊 Review Priority Levels

- `low` - Low priority, can wait
- `medium` - Normal priority (default)
- `high` - High priority, needs attention
- `critical` - Critical, urgent review needed

## 💡 Tips & Tricks

### View Pending Reviews Only
```bash
cool review histories --pending
```

### Quick Check Before Submitting to Collab
```bash
cool review submit-collab --list --pending
```

### Get Review ID from Histories
```bash
# Run this and copy the ID (first column)
cool review histories
```

### Cancel Submission
```bash
# When prompted "Submit this review to head architect? (yes/no):"
# Type: no
```

## 🔍 Common Workflows

### New Review Request Flow
```bash
1. cool review request
2. Fill in details interactively
3. Note the Request ID from output
4. Wait for tech lead approval
5. cool review submit-collab <request-id>
```

### Check Review Status
```bash
# See all reviews with status
cool review histories

# See only pending (not yet sent to architect)
cool review histories --pending

# See only completed (sent to architect)
cool review histories --completed
```

### Weekly Review Check
```bash
# See what needs to be submitted to collaboration
cool review histories --pending
```

## ❌ Common Errors & Solutions

### "missing user setup (name/email)"
```bash
# Solution: Run setup
cool setup email
```

### "missing GChat review webhook URL"
```bash
# Solution: Configure webhooks
cool setup webhook
```

### "missing GChat collaboration webhook URL"
```bash
# Solution: Configure webhooks
cool setup webhook
```

### "entry with ID xxx not found"
```bash
# Solution: Check ID from histories
cool review histories
# Copy the full ID or first 8 characters
```

### "review request already submitted to collaboration"
```bash
# This is expected - can't submit twice
# Check status with:
cool review histories
```

## 📝 Example Session

```bash
$ cool setup
🚀 Starting setup...
Name: Jane Developer
Email: jane@company.com
Review Webhook: https://chat.googleapis.com/v1/spaces/review/messages?key=xxx
Collab Webhook: https://chat.googleapis.com/v1/spaces/collab/messages?key=yyy
✅ Setup complete!

$ cool review request
Review Title: Fix authentication bug
Description: Resolve OAuth token expiration issue
Priority: critical
Pull Request Links:
  https://github.com/company/app/pull/456
  
Jira Ticket Links:
  https://jira.company.com/SEC-789
  
✅ Review request submitted successfully!
   Request ID: f3a2b1c4d5e6

$ cool review histories
ID        Title                       Priority   PRs  Jira  Submitted        Collab Status
--------  --------------------------  ---------  ---  ----  ---------------  -------------
f3a2b1c4  Fix authentication bug      critical   1    1     2025-01-15 14:30 ⏳ Pending

$ cool review submit-collab f3a2b1c4
📋 Review Request Details
Title: Fix authentication bug
Priority: critical
...
Submit this review to head architect? (yes/no): yes
✅ Successfully submitted to head architect!

$ cool review histories
ID        Title                       Priority   PRs  Jira  Submitted        Collab Status
--------  --------------------------  ---------  ---  ----  ---------------  -------------
f3a2b1c4  Fix authentication bug      critical   1    1     2025-01-15 14:30 ✅ Submitted
```

## 🎯 Keyboard Shortcuts

When filling review request:
- `Enter` - Confirm input
- `Ctrl+C` - Cancel operation
- Empty line - Finish list input (for PR/Jira links)

## 📱 Integration

### Google Chat Messages

**Review Request Message (to Tech Lead)**
```
🔍 *New Review Request*

*Title:* Fix authentication bug
*Priority:* critical
*Submitted by:* Jane Developer (jane@company.com)
*Submitted at:* 2025-01-15 14:30:00
...
```

**Collaboration Message (to Head Architect)**
```
🚀 *Review Request for Head Architect Approval*

*Title:* Fix authentication bug
*Priority:* critical
*Originally submitted by:* Jane Developer
*Tech Lead Approved:* ✅
*Forwarded at:* 2025-01-15 15:00:00
...
```

## 🔗 Links

- Full Spec: `docs/REVIEW_HISTORY_SPEC.md`
- Testing Guide: `docs/REVIEW_HISTORY_TESTING.md`
- Implementation Summary: `docs/IMPLEMENTATION_SUMMARY.md`

## 🆘 Help

For more details on any command:
```bash
cool --help
cool review --help
cool review request --help
cool review histories --help
cool review submit-collab --help
```

---

**Pro Tip**: Bookmark this guide for quick reference! 📌

