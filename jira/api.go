package jira

var defaultClient *JiraClient

func init() {
	defaultClient = NewJiraFromFile("jira.yaml")
}

func FetchBoards() []*Board {
	return defaultClient.FetchBoards()
}

func FetchSprints(boardId int) []*Sprint {
	return defaultClient.FetchSprints(boardId)
}

func FetchSprintDetails(boardId, sprintId int) *SprintDetails {
	return defaultClient.FetchSprintDetails(boardId, sprintId)
}

func GetBurndown(boardId, sprintId int, label string) *Burndown {
	sprintDetails := defaultClient.FetchSprintDetails(boardId, sprintId)
	issues := sprintDetails.Issues
	if label != "" {
		issues = Issues(issues).FilterByLabel(label)
	}

	return createBurndown(sprintDetails.Start, sprintDetails.End, issues)
}
