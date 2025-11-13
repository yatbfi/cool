# Preview Feature Implementation Summary

## ✅ Changes Completed

### 1. Priority System Updated (P0-P4)
- Changed from text input (`low/medium/high/critical`) to **select menu (1-5)**
- New priority levels:
  - **1 → P0 (Critical)**
  - **2 → P1 (High)**
  - **3 → P2 (Medium)** - Default
  - **4 → P3 (Low)**
  - **5 → P4 (Very Low)**

### 2. Preview Feature Added
- Added **preview before submit** functionality
- User can review the formatted message before sending
- Confirmation prompt: "Do you want to submit this review request? (yes/no)"
- Can cancel submission after preview

### 3. Usecase Updated with `withSend` Parameter
- `SubmitReviewRequest(ctx, req, withSend bool)` - Added `withSend` parameter
- When `withSend = false`: Returns preview entry without saving or sending
- When `withSend = true`: Saves to repository and sends to Google Chat
- Added `FormatReviewRequestMessage(entry)` method for formatting without duplication

## 📝 Updated Files

### 1. `internal/domain/usecase/review.go`
**Changes**:
- Updated `Review` interface with `withSend` parameter
- Added `FormatReviewRequestMessage()` method to interface
- Modified `SubmitReviewRequest()` to support preview mode
- Webhook validation only when `withSend = true`

**Key Code**:
```go
// Preview mode (withSend = false)
if !withSend {
    return entry, nil  // Returns without saving/sending
}

// Submit mode (withSend = true)
// Save to repository
if err := u.historyRepo.Save(ctx, entry); err != nil {
    return nil, fmt.Errorf("save history: %w", err)
}
// Send to GChat
message := formatReviewRequestMessage(entry)
if err := u.gchatUc.SendMessage(ctx, cfg.GChatReviewWebhookURL, message); err != nil {
    return nil, fmt.Errorf("send to GChat: %w", err)
}
```

### 2. `cmd/review_request.go`
**Changes**:
- Updated priority input to select menu (1-5)
- Added preview display before submission
- Added confirmation prompt
- Can cancel after preview

**New Flow**:
```
1. Collect input (title, description, priority, links)
2. Generate preview (withSend = false)
3. Display formatted preview message
4. Ask for confirmation
5. If yes → Submit (withSend = true)
   If no → Cancel
```

**Key Code**:
```go
// Preview first
previewEntry, err := c.reviewUc.SubmitReviewRequest(ctx, req, false)
message := c.reviewUc.FormatReviewRequestMessage(previewEntry)
fmt.Println(message)

// Confirmation
fmt.Print("Do you want to submit? (yes/no): ")
var confirmation string
fmt.Scanln(&confirmation)

if confirmation == "yes" || confirmation == "y" {
    // Submit for real
    entry, err := c.reviewUc.SubmitReviewRequest(ctx, req, true)
}
```

## 🎯 User Experience

### Before
```bash
$ cool review request
Review Title: Fix bug
Description: Fix critical bug
Priority (low/medium/high/critical) [medium]: high

⏳ Submitting review request...
✅ Review request submitted successfully!
```

### After
```bash
$ cool review request
Review Title: Fix bug
Description: Fix critical bug
Priority:
  1. P0 - Critical
  2. P1 - High
  3. P2 - Medium
  4. P3 - Low
  5. P4 - Very Low
Select priority [3]: 1

📋 Preview Review Request
=========================

🔍 *New Review Request*

*Title:* Fix bug
*Priority:* P0
*Submitted by:* John Doe (john@example.com)
*Submitted at:* 2025-11-13 10:30:00

*Description:*
Fix critical bug

*Request ID:* `abc123...`

Do you want to submit this review request? (yes/no): yes

⏳ Submitting review request...
✅ Review request submitted successfully!
```

## ✨ Benefits

### 1. No Duplicate Formatting
- Single `formatReviewRequestMessage()` function
- Used for both preview and sending
- No code duplication

### 2. Better UX
- User sees exactly what will be sent
- Can review and cancel before submitting
- Reduces mistakes and wrong submissions

### 3. Cleaner Priority Selection
- No typos (P0 vs p0 vs critical)
- Clear menu with descriptions
- Numbered selection is faster

### 4. Flexible Usecase
- `withSend` parameter makes usecase reusable
- Preview mode doesn't save or send
- Submit mode does everything
- Single method for both operations

## 🔧 Technical Details

### Method Signature Changes

**Before**:
```go
SubmitReviewRequest(ctx context.Context, req *ReviewRequest) (*ReviewHistoryEntry, error)
```

**After**:
```go
SubmitReviewRequest(ctx context.Context, req *ReviewRequest, withSend bool) (*ReviewHistoryEntry, error)
FormatReviewRequestMessage(entry *ReviewHistoryEntry) string
```

### Priority Mapping

**Before**:
- `low` → Low priority
- `medium` → Medium priority (default)
- `high` → High priority
- `critical` → Critical priority

**After**:
- `1` → `P0` (Critical)
- `2` → `P1` (High)
- `3` → `P2` (Medium) - default
- `4` → `P3` (Low)
- `5` → `P4` (Very Low)

## ✅ Build Status

- ✅ Build successful
- ✅ No errors
- ✅ No warnings
- ✅ All tests passing

## 🎊 Summary

**Preview feature dan P0-P4 priority system telah berhasil diimplementasikan!**

**Key Features**:
- ✅ Preview before submit
- ✅ Confirmation prompt
- ✅ P0-P4 priority selection menu
- ✅ No code duplication
- ✅ Cancel capability
- ✅ Better user experience

**Status**: Ready for use! 🚀

---

**Date**: November 13, 2025  
**Version**: 2.1.0 (Preview & P0-P4 Update)

