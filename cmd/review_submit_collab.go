package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/internal/domain/usecase"
	"github.com/yatbfi/cool/internal/pkg/table"
)

// ReviewSubmitCollabCmd handles submission to collaboration channel
type ReviewSubmitCollabCmd struct {
	*baseCmd
	reviewUc usecase.Review
	listOnly bool
	pending  bool
}

// NewReviewSubmitCollabCmd creates a new review submit-collab command
func NewReviewSubmitCollabCmd(reviewUc usecase.Review) *ReviewSubmitCollabCmd {
	cmd := &ReviewSubmitCollabCmd{
		reviewUc: reviewUc,
	}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "submit-collab [review-id]",
		Short: "Submit approved review to head architect",
		Long: `Submit an approved review request to head architect via collaboration channel.

Use this command after your review has been approved by tech lead.
The request will be forwarded to the head architect for final approval.

Examples:
  cool review submit-collab abc123            Submit specific review
  cool review submit-collab --list            Show all reviews with status
  cool review submit-collab --list --pending  Show only pending reviews`,
		RunE: cmd.run,
	})
	cmd.initFlags()
	return cmd
}

func (c *ReviewSubmitCollabCmd) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// If list flag is set, show history table
	if c.listOnly {
		return c.showHistoryTable(ctx)
	}

	// Validate arguments
	if len(args) == 0 {
		return fmt.Errorf("review ID is required\nUsage: cool review submit-collab <review-id>\nTip: Use --list to see available reviews")
	}

	reviewID := args[0]

	// Get review entry
	entry, err := c.reviewUc.GetHistoryByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("get review: %w", err)
	}

	// Display review details
	c.displayReviewDetails(entry)

	// Confirm submission
	if !c.confirmSubmission() {
		fmt.Println("\n❌ Submission cancelled")
		return nil
	}

	// Submit to collaboration
	fmt.Println("\n⏳ Submitting to collaboration channel...")
	if err := c.reviewUc.SubmitToCollaboration(ctx, reviewID); err != nil {
		return fmt.Errorf("submit to collaboration: %w", err)
	}

	// Success message
	fmt.Println()
	fmt.Println("✅ Successfully submitted to head architect!")
	fmt.Println()
	fmt.Println("💡 Your review request has been forwarded to the collaboration channel.")
	fmt.Println("   The head architect will review and provide approval.")
	fmt.Println()

	return nil
}

func (c *ReviewSubmitCollabCmd) showHistoryTable(ctx context.Context) error {
	filter := usecase.HistoryFilterAll
	if c.pending {
		filter = usecase.HistoryFilterPending
	}

	histories, err := c.reviewUc.GetHistories(ctx, filter)
	if err != nil {
		return fmt.Errorf("get histories: %w", err)
	}

	if len(histories) == 0 {
		fmt.Println()
		fmt.Println("📝 No review history found")
		if c.pending {
			fmt.Println("   No pending reviews to submit to collaboration.")
		}
		fmt.Println()
		return nil
	}

	c.displayHistoriesTable(histories)
	return nil
}

func (c *ReviewSubmitCollabCmd) displayHistoriesTable(histories []*usecase.ReviewHistoryEntry) {
	tbl := table.NewTable("ID", "Title", "Priority", "PRs", "Jira", "Submitted", "Collab Status")

	for _, entry := range histories {
		id := entry.ID
		if len(id) > 8 {
			id = id[:8]
		}

		title := entry.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		prCount := fmt.Sprintf("%d", len(entry.ReviewLinks))
		jiraCount := fmt.Sprintf("%d", len(entry.JiraLinks))
		submittedAt := entry.SubmittedAt.Format("2006-01-02 15:04")

		collabStatus := "⏳ Pending"
		if entry.SubmittedToCollab {
			collabStatus = "✅ Submitted"
		}

		tbl.AddRow(id, title, entry.Priority, prCount, jiraCount, submittedAt, collabStatus)
	}

	tbl.Print()
	fmt.Printf("Total: %d review(s)\n", tbl.RowCount())
	fmt.Println()
	fmt.Println("💡 To submit a review to collaboration: cool review submit-collab <id>")
	fmt.Println()
}

func (c *ReviewSubmitCollabCmd) displayReviewDetails(entry *usecase.ReviewHistoryEntry) {
	fmt.Println()
	fmt.Println("📋 Review Request Details")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Printf("ID: %s\n", entry.ID)
	fmt.Printf("Title: %s\n", entry.Title)
	fmt.Printf("Priority: %s\n", entry.Priority)
	fmt.Printf("Description: %s\n", entry.Description)
	fmt.Println()
	fmt.Printf("Submitted by: %s (%s)\n", entry.SubmittedBy, entry.SubmittedByEmail)
	fmt.Printf("Submitted at: %s\n", entry.SubmittedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	if len(entry.ReviewLinks) > 0 {
		fmt.Println("Pull Requests:")
		for _, link := range entry.ReviewLinks {
			fmt.Printf("  • %s\n", link)
		}
		fmt.Println()
	}

	if len(entry.JiraLinks) > 0 {
		fmt.Println("Jira Tickets:")
		for _, link := range entry.JiraLinks {
			fmt.Printf("  • %s\n", link)
		}
		fmt.Println()
	}

	if entry.SubmittedToCollab {
		fmt.Println("⚠️  This review has already been submitted to collaboration.")
		if entry.SubmittedToCollabAt != nil {
			fmt.Printf("   Submitted at: %s\n", entry.SubmittedToCollabAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	}
}

func (c *ReviewSubmitCollabCmd) confirmSubmission() bool {
	fmt.Println()
	fmt.Print("Submit this review to head architect? (yes/no): ")

	var response string
	_, _ = fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}

func (c *ReviewSubmitCollabCmd) initFlags() {
	flags := c.cmd.Flags()
	flags.BoolVarP(&c.listOnly, "list", "l", false, "Show review history table")
	flags.BoolVar(&c.pending, "pending", false, "Show only pending reviews (use with --list)")
}
