package jira

var defaultClient *JiraClient

func init() {
	defaultClient = NewJiraFromFile("jira.yaml")
}

func FetchBoardByName(boardName string) *Board {
	return defaultClient.FetchBoardByName(boardName)
}

func FetchSprint(boardName, sprintName string) *Sprint {
	return defaultClient.FetchSprint(boardName, sprintName)
}
