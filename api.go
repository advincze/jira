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
