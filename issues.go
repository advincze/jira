package jira

import "time"

type Issue struct {
	Id                      string
	Key                     string
	Type                    string
	OriginalEstimateSeconds int
	Labels                  []string
	Changes                 []*Change
}

type Issues []*Issue

type filterfunc func(*Issue) bool

func byType(isseType string) filterfunc {
	return func(issue *Issue) bool {
		return issue.Type == isseType
	}
}

func byLabel(labelToSearch string) filterfunc {
	return func(issue *Issue) (containsLabel bool) {
		for _, labelFound := range issue.Labels {
			if labelFound == labelToSearch {
				containsLabel = true
				break
			}
		}
		return
	}
}

func (issues Issues) Filter(fn filterfunc) Issues {
	filteredIssues := make([]*Issue, 0, len(issues))
	for _, issue := range issues {
		if fn(issue) {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return Issues(filteredIssues)
}

func (issues Issues) FilterByType(issueType string) Issues {
	return issues.Filter(byType(issueType))
}

func (issues Issues) FilterByLabel(labelToSearch string) Issues {
	return issues.Filter(byLabel(labelToSearch))
}

func (issues Issues) filterByLabel2(labelToSearch string) Issues {
	filteredIssues := make([]*Issue, 0, len(issues))
	for _, issue := range issues {
		var containsLabel bool
		for _, labelFound := range issue.Labels {
			if labelFound == labelToSearch {
				containsLabel = true
				break
			}
		}
		if containsLabel {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return Issues(filteredIssues)
}

type Change struct {
	Timestamp time.Time
	Field     string
	From      string
	To        string
}

func (c *Change) issueClosed() bool {

	if c.Field == "status" {
		switch c.From {
		case "Open":
			break
		case "Planung":
			break
		case "In Progress":
			break
		case "Geschlossen":
			return false
		}

		switch c.To {
		case "Open":
			break
		case "Planung":
			break
		case "In Progress":
			break
		case "Geschlossen":

			return true
		}

	}
	return false
}
