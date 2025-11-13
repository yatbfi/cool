# Review History System - Testing Guide

## Build & Installation

### Build from Source

```bash
cd /Users/mac-209055/go/src/github.com/yatbfi/cool
go mod tidy
go build -o bin/cool .
```

### Install Globally

```bash
go install github.com/yatbfi/cool@latest
```

## Testing Commands

### 1. Setup Commands

#### Setup Email

```bash
./bin/cool setup email
```

**Expected**:
- Prompt for name
- Prompt for email
- Save to `~/.cool-cli/config.json`

#### Setup Webhook

```bash
./bin/cool setup webhook
```

**Expected**:
- Prompt for GChat Review Webhook URL
- Prompt for GChat Collab Webhook URL
- Save to `~/.cool-cli/config.json`

#### Full Setup

```bash
./bin/cool setup
```

**Expected**:
- Run email setup if not configured
- Run webhook setup if not configured
- Skip if already configured

### 2. Review Commands

#### Submit Review Request

```bash
./bin/cool review request
```

**Expected Flow**:
1. Prompt for review title (required)
2. Prompt for description (optional)
3. Prompt for priority (low/medium/high/critical, default: medium)
4. Prompt for PR links (multiple, empty line to finish)
5. Prompt for Jira links (multiple, empty line to finish)
6. Generate unique ID
7. Save to `~/.cool-cli/review_histories.json`
8. Send message to GChat review webhook
9. Display success with ID

**Test Cases**:
- ✅ Valid submission with all fields
- ✅ Valid submission with minimal fields (title + priority)
- ❌ Missing title (should error)
- ❌ Invalid priority (should error)
- ❌ Webhook not configured (should error with helpful message)

#### View All Histories

```bash
./bin/cool review histories
```

**Expected**:
- Display table with all review histories
- Columns: ID, Title, Priority, PRs, Jira, Submitted, Collab Status, Collab Submitted
- Show total count
- Show helpful hint for next action

**Test Cases**:
- ✅ Display multiple entries
- ✅ Display empty message if no histories
- ✅ Truncate long titles (>40 chars)
- ✅ Truncate long IDs (show first 8 chars)
- ✅ Format timestamps correctly

#### View Pending Histories

```bash
./bin/cool review histories --pending
```

**Expected**:
- Display only reviews not yet submitted to collaboration
- Same table format as above
- Show count of pending reviews

**Test Cases**:
- ✅ Filter correctly by collab status
- ✅ Show empty message if no pending reviews

#### View Completed Histories

```bash
./bin/cool review histories --completed
```

**Expected**:
- Display only reviews already submitted to collaboration
- Show collaboration submission timestamps
- Same table format as above

**Test Cases**:
- ✅ Filter correctly by collab status
- ✅ Show empty message if no completed reviews

#### Submit to Collaboration

```bash
./bin/cool review submit-collab <review-id>
```

**Expected Flow**:
1. Validate collab webhook configured
2. Load review entry by ID
3. Check not already submitted
4. Display review details
5. Ask for confirmation
6. Send message to GChat collab webhook
7. Update history entry (set submitted flag & timestamp)
8. Display success message

**Test Cases**:
- ✅ Valid submission with confirmation "yes"
- ✅ Valid submission with confirmation "y"
- ❌ Cancel with "no" or "n"
- ❌ Invalid review ID (should error)
- ❌ Already submitted (should error with message)
- ❌ Webhook not configured (should error with helpful message)

#### List with Submit-Collab

```bash
./bin/cool review submit-collab --list
./bin/cool review submit-collab --list --pending
```

**Expected**:
- Same as `review histories` commands
- Alternative way to view histories
- Helpful for discovering review IDs to submit

## Manual Testing Scenarios

### Scenario 1: Complete Workflow

```bash
# 1. Initial setup
./bin/cool setup email
# Enter: John Doe, john@example.com

./bin/cool setup webhook
# Enter review webhook: https://chat.googleapis.com/v1/spaces/.../messages?key=...
# Enter collab webhook: https://chat.googleapis.com/v1/spaces/.../messages?key=...

# 2. Submit first review
./bin/cool review request
# Title: Feature X Implementation
# Description: Add new feature X with tests
# Priority: high
# PR: https://github.com/org/repo/pull/123
# PR: (empty line)
# Jira: https://jira.example.com/PROJ-456
# Jira: (empty line)

# 3. View histories
./bin/cool review histories
# Should show 1 entry with "Pending" status

# 4. Submit second review
./bin/cool review request
# Title: Bug Fix for Y
# Description: Fix critical bug in Y
# Priority: critical
# PR: https://github.com/org/repo/pull/124
# Jira: https://jira.example.com/PROJ-457

# 5. View pending only
./bin/cool review histories --pending
# Should show 2 entries

# 6. Submit to collaboration
./bin/cool review submit-collab abc123...
# Confirm: yes

# 7. View histories again
./bin/cool review histories
# Should show 2 entries, 1 with "Submitted" status

# 8. View completed
./bin/cool review histories --completed
# Should show 1 entry

# 9. Try to submit again (should fail)
./bin/cool review submit-collab abc123...
# Should error: already submitted
```

### Scenario 2: Error Handling

```bash
# Test without setup
./bin/cool review request
# Should error: missing user setup

./bin/cool setup email
# Complete email setup

./bin/cool review request
# Should error: missing webhook URL

./bin/cool setup webhook
# Complete webhook setup

./bin/cool review request
# Should now work
```

### Scenario 3: Validation

```bash
# Test invalid priority
./bin/cool review request
# Title: Test
# Priority: invalid
# Should error: invalid priority

# Test empty title
./bin/cool review request
# Title: (empty)
# Should error: title is required

# Test invalid review ID
./bin/cool review submit-collab invalid-id
# Should error: not found
```

## File Verification

### Config File

```bash
cat ~/.cool-cli/config.json
```

**Expected Structure**:
```json
{
  "gchat_review_webhook_url": "https://...",
  "gchat_collab_webhook_url": "https://...",
  "user_name": "John Doe",
  "user_email": "john@example.com"
}
```

### History File

```bash
cat ~/.cool-cli/review_histories.json
```

**Expected Structure**:
```json
[
  {
    "id": "abc123...",
    "title": "Feature X Implementation",
    "description": "Add new feature X with tests",
    "priority": "high",
    "review_links": ["https://github.com/org/repo/pull/123"],
    "jira_links": ["https://jira.example.com/PROJ-456"],
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

## Integration Testing with Google Chat

### Setup Test Webhooks

1. Create two Google Chat spaces or use existing ones
2. Create incoming webhooks for each space
3. Configure them using `cool setup webhook`

### Test Review Webhook

```bash
./bin/cool review request
# Fill in details
# Check Google Chat space for message
```

**Expected Message Format**:
```
🔍 *New Review Request*

*Title:* Feature X Implementation
*Priority:* high
*Submitted by:* John Doe (john@example.com)
*Submitted at:* 2025-01-15 10:30:00

*Description:*
Add new feature X with tests

*Review Links:*
• https://github.com/org/repo/pull/123

*Jira Links:*
• https://jira.example.com/PROJ-456

*Request ID:* `abc123...`
```

### Test Collab Webhook

```bash
./bin/cool review submit-collab abc123...
# Confirm: yes
# Check Google Chat space for message
```

**Expected Message Format**:
```
🚀 *Review Request for Head Architect Approval*

*Title:* Feature X Implementation
*Priority:* high
*Originally submitted by:* John Doe (john@example.com)
*Tech Lead Approved:* ✅
*Forwarded at:* 2025-01-15 14:00:00

*Description:*
Add new feature X with tests

*Review Links:*
• https://github.com/org/repo/pull/123

*Jira Links:*
• https://jira.example.com/PROJ-456

*Request ID:* `abc123...`

_Please review and approve._
```

## Performance Testing

### Test with Large History

```bash
# Create script to generate multiple entries
for i in {1..50}; do
  echo -e "Review $i\nDescription $i\nhigh\nhttps://github.com/pull/$i\n\nhttps://jira.com/PROJ-$i\n" | ./bin/cool review request
done

# View histories
./bin/cool review histories
# Should display all 50 entries with good performance
```

### Concurrent Access

```bash
# Test concurrent submissions (in different terminals)
Terminal 1: ./bin/cool review request
Terminal 2: ./bin/cool review request
Terminal 3: ./bin/cool review histories
# Should handle concurrent access safely
```

## Regression Testing Checklist

- [ ] Build succeeds without errors
- [ ] All commands show help correctly
- [ ] Setup commands save configuration properly
- [ ] Review request validates input correctly
- [ ] Review request sends to webhook successfully
- [ ] History is saved correctly
- [ ] Histories display correctly with all formats
- [ ] Filtering works (--pending, --completed)
- [ ] Submit-collab validates webhook configuration
- [ ] Submit-collab prevents duplicate submission
- [ ] Submit-collab sends to webhook successfully
- [ ] Submit-collab updates history correctly
- [ ] Error messages are helpful and actionable
- [ ] File permissions are correct (0644 for files, 0755 for dirs)
- [ ] Concurrent access doesn't corrupt data
- [ ] Large histories display without issues

## Known Limitations

1. No edit/delete functionality for histories (future enhancement)
2. No approval tracking integration (future enhancement)
3. No search functionality (future enhancement)
4. IDs are not human-readable (could use shorter format)
5. No pagination for large histories

## Troubleshooting

### Issue: Build Fails

```bash
# Check Go version
go version  # Should be 1.21+

# Clean and rebuild
go clean
go mod tidy
go build -o bin/cool .
```

### Issue: Permission Denied

```bash
chmod +x bin/cool
```

### Issue: Config Not Found

```bash
# Check if directory exists
ls -la ~/.cool-cli/

# Run setup
./bin/cool setup
```

### Issue: Webhook Not Working

```bash
# Test webhook manually
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"text": "Test message"}' \
  'YOUR_WEBHOOK_URL'
```

### Issue: History File Corrupted

```bash
# Backup and recreate
mv ~/.cool-cli/review_histories.json ~/.cool-cli/review_histories.json.bak
echo '[]' > ~/.cool-cli/review_histories.json
```

## Success Criteria

✅ All commands execute without crashes
✅ All validations work as expected
✅ Error messages are clear and helpful
✅ Data is persisted correctly
✅ Webhooks receive properly formatted messages
✅ Concurrent access is safe
✅ Performance is acceptable for typical usage (100+ entries)

