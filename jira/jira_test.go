package jira

import (
	"testing"
)

func TestFetchRapidViews(t *testing.T) {

	rapidViews := defaultClient.FetchViews()

	t.Logf("fetched %d views. \n", len(rapidViews.Views))
	if len(rapidViews.Views) == 0 {
		t.Error("rapidViews should not be empty")
	}
}

func TestFetchSprints(t *testing.T) {

	boardId := 51
	sprints := defaultClient.FetchSprints(boardId)

	t.Logf("fetched %d sprints for boardId: %d \n", len(sprints.Sprints), boardId)
	// t.Logf("sprints: %q \n", sprints)
	for _, sprint := range sprints.Sprints {
		t.Logf("sprint: %v, %v %v \n", sprint.Id, sprint.Name)
	}

	if len(sprints.Sprints) == 0 {
		t.Error("sprints should not be empty")
	}
}

func TestFetchSprintDetails(t *testing.T) {
	boardId := 51
	sprintId := 217

	sprintDetails := defaultClient.FetchSprintDetails(boardId, sprintId)
	s := sprintDetails.Sprint
	t.Logf("sprint details: %q \n\n", s)
	if s.Id != 217 {
		t.Errorf("sprint id was %d shoud be %d ", s.Id, sprintId)
	}
	// sprintDetails.Sprint.
}

// func _TestSearch(t *testing.T) {

// 	sprintDetails := defaultClient.FetchSprintDetails(51, 217)

// 	keys := make([]string, 0, 100)
// 	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
// 		keys = append(keys, Issue.Key)
// 	}
// 	issues := defaultClient.FetchIssues(keys)

// 	for _, issue := range issues.Issues {
// 		for _, item := range issue.Changelog.Histories.Items {
// 			fmt.Println(issue.Key, issue.Fields.Summary, item.Field, item.FromString, item.ToString)
// 		}
// 	}

// }
