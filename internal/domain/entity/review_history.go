package entity

import "time"

// ReviewHistoryEntry represents a single review request history entry
type ReviewHistoryEntry struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	Priority            string     `json:"priority"`
	ReviewLinks         []string   `json:"review_links"`
	JiraLinks           []string   `json:"jira_links"`
	SubmittedBy         string     `json:"submitted_by"`
	SubmittedByEmail    string     `json:"submitted_by_email"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	SubmittedToCollab   bool       `json:"submitted_to_collab"`
	SubmittedToCollabAt *time.Time `json:"submitted_to_collab_at,omitempty"`
	SubmittedToCollabBy string     `json:"submitted_to_collab_by,omitempty"`
	ApprovedByTechLead  bool       `json:"approved_by_tech_lead"`
	ApprovedByArchitect bool       `json:"approved_by_architect"`
	Notes               string     `json:"notes,omitempty"`
}
