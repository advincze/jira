package jira

import "testing"

func TestInit(t *testing.T) {

	SetConfig(&Config{})

	jiraGetRequestFetcher = func(*JiraClient, string) []byte {
		return []byte("{}")
	}
	boards := FetchBoards()
	if len(boards) != 0 {
		t.Errorf("should be empty")
	}
}
