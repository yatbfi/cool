package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/internal/domain/usecase"
	"github.com/yatbfi/cool/internal/pkg/table"
)

// ReviewHistoryCmd displays review history table
type ReviewHistoryCmd struct {
	*baseCmd
	reviewUc  usecase.Review
	pending   bool
	completed bool
}

// NewReviewHistoriesCmd creates a new review histories command
func NewReviewHistoriesCmd(reviewUc usecase.Review) *ReviewHistoryCmd {
	cmd := &ReviewHistoryCmd{
		reviewUc: reviewUc,
	}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "history",
		Short: "Display review request history",
		Long: `Display all review request history with their submission status.

Use flags to filter the results:
- --pending: Show only reviews not yet submitted to collaboration
- --completed: Show only reviews already submitted to collaboration

Examples:
  cool review histories              # Show all histories
  cool review histories --pending    # Show pending only
  cool review histories --completed  # Show completed only`,
		RunE: cmd.run,
	})
	cmd.initFlags()
	return cmd
}

func (c *ReviewHistoryCmd) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	// Determine filter
	filter := usecase.HistoryFilterAll
	if c.pending && c.completed {
		return fmt.Errorf("cannot use both --pending and --completed flags")
	}
	if c.pending {
		filter = usecase.HistoryFilterPending
	} else if c.completed {
		filter = usecase.HistoryFilterCompleted
	}

	// Get histories
	histories, err := c.reviewUc.GetHistories(ctx, filter)
	if err != nil {
		return fmt.Errorf("get histories: %w", err)
	}

	// Display
	if len(histories) == 0 {
		c.displayEmptyMessage(filter)
	} else {
		c.displayHistoriesTable(histories)
	}

	return nil
}

func (c *ReviewHistoryCmd) displayEmptyMessage(filter usecase.HistoryFilter) {
	fmt.Println()
	fmt.Println("📝 No review history found")
	switch filter {
	case usecase.HistoryFilterPending:
		fmt.Println("   No pending reviews to submit to collaboration.")
	case usecase.HistoryFilterCompleted:
		fmt.Println("   No completed reviews submitted to collaboration.")
	default:
		fmt.Println("   Submit your first review request using:")
		fmt.Println("   cool review request")
	}
	fmt.Println()
}

func (c *ReviewHistoryCmd) displayHistoriesTable(histories []*usecase.ReviewHistoryEntry) {
	tbl := table.NewTable("ID", "Title", "Priority", "PRs", "Jira", "Submitted", "Collab Status", "Collab Submitted")

	for _, entry := range histories {
		id := entry.ID
		title := entry.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		prCount := fmt.Sprintf("%d", len(entry.ReviewLinks))
		jiraCount := fmt.Sprintf("%d", len(entry.JiraLinks))
		submittedAt := entry.SubmittedAt.Format("2006-01-02 15:04")

		collabStatus := "⏳ Pending"
		collabSubmitted := "-"
		if entry.SubmittedToCollab {
			collabStatus = "✅ Submitted"
			if entry.SubmittedToCollabAt != nil {
				collabSubmitted = entry.SubmittedToCollabAt.Format("2006-01-02 15:04")
			}
		}

		tbl.AddRow(id, title, entry.Priority, prCount, jiraCount, submittedAt, collabStatus, collabSubmitted)
	}

	tbl.Print()
	fmt.Printf("Total: %d review(s)\n", tbl.RowCount())
	fmt.Println()
	fmt.Println("💡 To view details or submit to collaboration: cool review submit-collab <id>")
	fmt.Println()
}

func (c *ReviewHistoryCmd) initFlags() {
	flags := c.cmd.Flags()
	flags.BoolVar(&c.pending, "pending", false, "Show only pending reviews (not submitted to collaboration)")
	flags.BoolVar(&c.completed, "completed", false, "Show only completed reviews (submitted to collaboration)")
}
