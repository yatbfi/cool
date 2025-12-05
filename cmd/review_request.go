package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/config"
	"github.com/yatbfi/cool/internal/domain/usecase"
	"github.com/yatbfi/cool/internal/pkg/common"
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

	// Collect input first time
	req, err := c.collectReviewRequest()
	if err != nil {
		return err
	}

	// Main loop for preview and confirmation
	for {
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

		// Ask for confirmation with edit option
		fmt.Print("Do you want to (s)ubmit, (e)dit, or (c)ancel? [s/e/c]: ")
		var action string
		_, _ = fmt.Scanln(&action)
		action = strings.ToLower(strings.TrimSpace(action))

		switch action {
		case "s", "submit", "":
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

		case "e", "edit":
			// Ask what to edit
			editReq, err := c.editReviewRequest(req)
			if err != nil {
				return err
			}
			req = editReq
			// Continue loop to show preview again

		case "c", "cancel":
			fmt.Println("\n❌ Review request cancelled")
			return nil

		default:
			fmt.Println("❌ Invalid option. Please choose (s)ubmit, (e)dit, or (c)ancel.")
			fmt.Println()
			// Continue loop
		}
	}
}

func (c *ReviewRequestCmd) collectReviewRequest() (*usecase.ReviewRequest, error) {
	reader := bufio.NewReader(os.Stdin)

	// Title - keep as single line input for simplicity
	fmt.Print("Review Title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read title: %w", err)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Description - use editor for multiline input
	fmt.Println()

	// Get editor command from config
	cfg := config.GetConfig()
	editorCmd, err := common.GetEditorCommand(cfg.PreferredEditor)
	if err != nil {
		return nil, fmt.Errorf("get editor: %w. Run 'cool setup editor' to configure", err)
	}

	fmt.Printf("Opening editor (%s) for description...\n", common.GetEditorDisplayName(editorCmd))
	description, err := common.OpenEditor(editorCmd, "Enter your review description below")
	if err != nil {
		return nil, fmt.Errorf("open editor: %w", err)
	}

	fmt.Println("✓ Description captured")

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

func (c *ReviewRequestCmd) editReviewRequest(req *usecase.ReviewRequest) (*usecase.ReviewRequest, error) {
	reader := bufio.NewReader(os.Stdin)
	cfg := config.GetConfig()

	fmt.Println()
	fmt.Println("What would you like to edit?")
	fmt.Println("  1. Title")
	fmt.Println("  2. Description")
	fmt.Println("  3. Priority")
	fmt.Println("  4. Pull Request Links")
	fmt.Println("  5. Jira Ticket Links")
	fmt.Print("Select field to edit [1-5]: ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return req, fmt.Errorf("read choice: %w", err)
	}
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		// Edit Title
		fmt.Printf("Current title: %s\n", req.Title)
		fmt.Print("New title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			return req, fmt.Errorf("read title: %w", err)
		}
		title = strings.TrimSpace(title)
		if title != "" {
			req.Title = title
		}

	case "2":
		// Edit Description
		editorCmd, err := common.GetEditorCommand(cfg.PreferredEditor)
		if err != nil {
			return req, fmt.Errorf("get editor: %w", err)
		}
		fmt.Printf("\nOpening editor (%s) for description...\n", common.GetEditorDisplayName(editorCmd))
		description, err := common.OpenEditorWithContent(editorCmd, req.Description)
		if err != nil {
			return req, fmt.Errorf("open editor: %w", err)
		}
		req.Description = description
		fmt.Println("✓ Description updated")

	case "3":
		// Edit Priority
		fmt.Printf("Current priority: %s\n", req.Priority)
		for {
			fmt.Println("Priority:")
			fmt.Println("  1. P0 - Critical")
			fmt.Println("  2. P1 - High")
			fmt.Println("  3. P2 - Medium")
			fmt.Println("  4. P3 - Low")
			fmt.Println("  5. P4 - Very Low")
			fmt.Print("Select priority: ")
			priorityInput, err := reader.ReadString('\n')
			if err != nil {
				return req, fmt.Errorf("read priority: %w", err)
			}
			priorityInput = strings.TrimSpace(priorityInput)

			priority, err := getPriorityFromSelection(priorityInput)
			if err != nil {
				fmt.Printf("❌ %s. Please try again.\n\n", err.Error())
				continue
			}
			req.Priority = priority
			break
		}

	case "4":
		// Edit PR Links
		fmt.Println("\nCurrent Pull Request Links:")
		for i, link := range req.ReviewLinks {
			fmt.Printf("  %d. %s\n", i+1, link)
		}
		fmt.Println("\nEnter new Pull Request Links (one per line, empty line to finish):")
		req.ReviewLinks = c.collectLinks(reader)

	case "5":
		// Edit Jira Links
		fmt.Println("\nCurrent Jira Ticket Links:")
		for i, link := range req.JiraLinks {
			fmt.Printf("  %d. %s\n", i+1, link)
		}
		fmt.Println("\nEnter new Jira Ticket Links (one per line, empty line to finish):")
		req.JiraLinks = c.collectLinks(reader)

	default:
		fmt.Println("❌ Invalid choice. No changes made.")
	}

	return req, nil
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
