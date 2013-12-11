package jira

import (
	"testing"
)

var testClient *JiraClient

func init() {
	testClient = NewJiraFromFile("../jira.yaml")
}

func TestFetchRapidViews(t *testing.T) {
	rapidViews := testClient.fetchRapidViews()

	t.Logf("fetched %d views. \n", len(rapidViews.Views))
	if len(rapidViews.Views) == 0 {
		t.Error("rapidViews should not be empty")
	}
}

// func TestFetchSprints(t *testing.T) {

// 	boardId := 51
// 	sprints := defaultClient.FetchSprints(boardId)

// 	t.Logf("fetched %d sprints for boardId: %d \n", len(sprints.Sprints), boardId)
// 	// t.Logf("sprints: %q \n", sprints)
// 	for _, sprint := range sprints.Sprints {
// 		t.Logf("sprint: %v, %v \n", sprint.Id, sprint.Name)
// 	}

// 	if len(sprints.Sprints) == 0 {
// 		t.Error("sprints should not be empty")
// 	}
// }

// func TestFetchSprintDetails(t *testing.T) {
// 	boardId := 51
// 	sprintId := 217
// 	sprintDetails := defaultClient.FetchSprintDetails(boardId, sprintId)
// 	s := sprintDetails.Sprint
// 	t.Logf("sprint details: %q \n\n", s)
// 	if s.Id != 217 {
// 		t.Errorf("sprint id was %d shoud be %d ", s.Id, sprintId)
// 	}
// }

// func TestFetchSprintIssues(t *testing.T) {
// 	sprintId := 217
// 	sprintIssues := defaultClient.FetchSprintIssues(sprintId)
// 	t.Logf("sprint issues: %q \n\n", sprintIssues)
// }
