package common

import "strings"

func ExtractJiraTicketNumber(url string) string {
	parts := strings.Split(url, "/browse/")
	if len(parts) < 2 {
		return ""
	}
	ticket := strings.TrimSpace(parts[1])
	if idx := strings.IndexAny(ticket, "?#"); idx > 0 {
		ticket = ticket[:idx]
	}
	return ticket
}

func ExtractPRNumber(url string) string {
	for _, segment := range []string{"/pull/", "/merge_requests/"} { // handles GitHub and GitLab
		if strings.Contains(url, segment) {
			parts := strings.Split(url, segment)
			if len(parts) < 2 {
				continue
			}
			num := parts[1]
			if idx := strings.IndexAny(num, "?#"); idx > 0 {
				num = num[:idx]
			}
			return strings.Trim(num, "/")
		}
	}
	return ""
}
