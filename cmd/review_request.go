package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/internal/domain/usecase"
)

// ReviewRequestCmd handles review request submission
type ReviewRequestCmd struct {
	*baseCmd
	reviewUc usecase.Review
}

// NewReviewRequestCmd creates a new review request command
func NewReviewRequestCmd(reviewUc usecase.Review) *ReviewRequestCmd {
	cmd := &ReviewRequestCmd{
		reviewUc: reviewUc,
	}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "request",
		Short: "Submit a new review request to tech lead",
		Long: `Submit a new code review request to tech lead.

This command will prompt you for:
- Review title
- Description
- Priority (P0-P4: Critical to Very Low)
- Pull request links
- Jira ticket links

The request will be saved in your review history and sent to the configured Google Chat webhook.`,
		RunE: cmd.run,
	})
	return cmd
}

func (c *ReviewRequestCmd) run(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	fmt.Println()
	fmt.Println("📝 Submit Review Request to Tech Lead")
	fmt.Println("=====================================")
	fmt.Println()

	// Collect input
	req, err := c.collectReviewRequest()
	if err != nil {
		return err
	}

	// Preview first (without sending)
	fmt.Println()
	fmt.Println("📋 Preview Review Request")
	fmt.Println("=========================")

	previewEntry, err := c.reviewUc.SubmitReviewRequest(ctx, req, false)
	if err != nil {
		return fmt.Errorf("generate preview: %w", err)
	}

	// Display preview message
	message := c.reviewUc.FormatReviewRequestMessage(previewEntry)
	fmt.Println()
	fmt.Println(message)
	fmt.Println()

	// Ask for confirmation
	fmt.Print("Do you want to submit this review request? (yes/no): ")
	var confirmation string
	_, _ = fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(strings.TrimSpace(confirmation))

	if confirmation != "yes" && confirmation != "y" {
		fmt.Println("\n❌ Review request cancelled")
		return nil
	}

	// Submit request
	fmt.Println()
	fmt.Println("⏳ Submitting review request...")

	entry, err := c.reviewUc.SubmitReviewRequest(ctx, req, true)
	if err != nil {
		return fmt.Errorf("submit review request: %w", err)
	}

	// Display success
	fmt.Println()
	fmt.Println("✅ Review request submitted successfully!")
	fmt.Println()
	fmt.Printf("   Request ID: %s\n", entry.ID)
	fmt.Printf("   Title: %s\n", entry.Title)
	fmt.Printf("   Priority: %s\n", entry.Priority)
	fmt.Printf("   Submitted at: %s\n", entry.SubmittedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("💡 Your request has been sent to tech lead for review.")
	fmt.Println("   Once approved, you can forward it to head architect using:")
	fmt.Printf("   cool review submit-collab %s\n", entry.ID)
	fmt.Println()

	return nil
}

func (c *ReviewRequestCmd) collectReviewRequest() (*usecase.ReviewRequest, error) {
	reader := bufio.NewReader(os.Stdin)

	// Title
	fmt.Print("Review Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read title: %w", err)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Description
	fmt.Print("Description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read description: %w", err)
	}
	description = strings.TrimSpace(description)

	// Priority
	var priority string
	for {
		fmt.Println("Priority:")
		fmt.Println("  1. P0 - Critical")
		fmt.Println("  2. P1 - High")
		fmt.Println("  3. P2 - Medium")
		fmt.Println("  4. P3 - Low")
		fmt.Println("  5. P4 - Very Low")
		fmt.Print("Select priority [3]: ")
		priorityInput, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("read priority: %w", err)
		}
		priorityInput = strings.TrimSpace(priorityInput)
		if priorityInput == "" {
			priorityInput = "3"
		}

		priority, err = getPriorityFromSelection(priorityInput)
		if err != nil {
			fmt.Printf("❌ %s. Please try again.\n\n", err.Error())
			continue
		}
		break
	}

	// Review Links
	fmt.Println("Pull Request Links (one per line, empty line to finish):")
	reviewLinks := c.collectLinks(reader)

	// Jira Links
	fmt.Println("Jira Ticket Links (one per line, empty line to finish):")
	jiraLinks := c.collectLinks(reader)

	return &usecase.ReviewRequest{
		Title:       title,
		Description: description,
		Priority:    priority,
		ReviewLinks: reviewLinks,
		JiraLinks:   jiraLinks,
	}, nil
}

func (c *ReviewRequestCmd) collectLinks(reader *bufio.Reader) []string {
	var links []string
	for {
		fmt.Print("  ")
		link, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		link = strings.TrimSpace(link)
		if link == "" {
			break
		}
		links = append(links, link)
	}
	return links
}

func getPriorityFromSelection(selection string) (string, error) {
	priorities := map[string]string{
		"1": "P0",
		"2": "P1",
		"3": "P2",
		"4": "P3",
		"5": "P4",
	}

	priority, ok := priorities[selection]
	if !ok {
		return "", fmt.Errorf("invalid priority selection: %s (must be 1-5)", selection)
	}

	return priority, nil
}
