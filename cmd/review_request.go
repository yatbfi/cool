package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yatbfi/cool/internal/domain/usecases"
)

type ReviewRequestCmd struct {
	*baseCmd
	reviewUc usecases.Review
}

func NewReviewRequestCmd(
	reviewUc usecases.Review,
) *ReviewRequestCmd {
	cmd := &ReviewRequestCmd{
		reviewUc: reviewUc,
	}
	cmd.baseCmd = newBaseCommand(&cobra.Command{
		Use:   "request",
		Short: "Submit request review to channel code review",
		RunE:  cmd.run,
	})
	return cmd
}

func (c *ReviewRequestCmd) run(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// 1. Title
	titlePrompt := promptui.Prompt{
		Label: "Changes title",
	}
	title, err := titlePrompt.Run()
	if err != nil {
		return err
	}

	// 2. PR Links
	var reviewLinks []*usecases.ReviewLink
	for {
		addPrompt := promptui.Prompt{
			Label:     "Add PR link?",
			IsConfirm: true,
		}
		_, err := addPrompt.Run()
		if err != nil {
			break
		}

		servicePrompt := promptui.Prompt{
			Label: "PR Service (e.g., LSS, LGS, LTS, LTW, LPW)",
		}
		service, err := servicePrompt.Run()
		if err != nil {
			return err
		}

		linkPrompt := promptui.Prompt{
			Label: "PR Link",
		}
		link, err := linkPrompt.Run()
		if err != nil {
			return err
		}

		reviewLinks = append(reviewLinks, &usecases.ReviewLink{
			Service:        service,
			PullRequestURL: link,
		})
	}

	// 3. Jira Links
	var jiraLinks []string
	for {
		addPrompt := promptui.Prompt{
			Label:     "Add Jira link?",
			IsConfirm: true,
		}
		_, err := addPrompt.Run()
		if err != nil {
			break
		}

		linkPrompt := promptui.Prompt{
			Label: "Jira Link",
		}
		link, err := linkPrompt.Run()
		if err != nil {
			return err
		}

		jiraLinks = append(jiraLinks, link)
	}

	// 4. Description
	descPrompt := promptui.Prompt{
		Label: "Description",
	}
	description, err := descPrompt.Run()
	if err != nil {
		return err
	}

	// 5. Priority (select)
	priorityPrompt := promptui.Select{
		Label: "Priority",
		Items: []string{"P0", "P1", "P2", "P3"},
	}
	_, priority, err := priorityPrompt.Run()
	if err != nil {
		return err
	}

	return c.reviewUc.RequestReview(ctx, &usecases.RequestReviewPayload{
		Title:       fmt.Sprintf("[%s] %s", priority, title),
		Description: description,
		ReviewLinks: reviewLinks,
		JiraLinks:   jiraLinks,
	})
}
