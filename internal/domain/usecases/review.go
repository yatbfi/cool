package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/yatbfi/cool/config"
)

type RequestReviewPayload struct {
	Title       string
	Description string
	ReviewLinks []*ReviewLink
	JiraLinks   []string
}

type ReviewLink struct {
	Service        string
	PullRequestURL string
}

type Review interface {
	RequestReview(ctx context.Context, p *RequestReviewPayload) error
}

type ReviewUsecase struct {
	cfg   *config.Config
	gChat GChat
}

var _ Review = (*ReviewUsecase)(nil)

func NewReviewUsecase(
	cfg *config.Config,
	gChat GChat,
) *ReviewUsecase {
	return &ReviewUsecase{
		cfg:   cfg,
		gChat: gChat,
	}
}

func (u *ReviewUsecase) RequestReview(ctx context.Context, p *RequestReviewPayload) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*%s*\n\n", p.Title))
	sb.WriteString(fmt.Sprintf("_%s_\n\n", p.Description))

	sb.WriteString("*Pull Requests:*\n")
	for _, rl := range p.ReviewLinks {
		sb.WriteString(fmt.Sprintf("• %s → <%s|Open>\n", rl.Service, rl.PullRequestURL))
	}
	sb.WriteString("\n")

	if len(p.JiraLinks) > 0 {
		sb.WriteString("*Jira Links:*\n")
		for _, link := range p.JiraLinks {
			sb.WriteString(fmt.Sprintf("• <%s|Jira Ticket>\n", link))
		}
		sb.WriteString("\n")
	}

	message := sb.String()

	if err := u.gChat.SendToChannel(ctx, "lora_code_review", message); err != nil {
		return fmt.Errorf("failed to send review message: %w", err)
	}

	fmt.Println("✅ Review request sent successfully!")
	return nil
}
