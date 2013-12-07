package jira

import (
	"fmt"
	"testing"
)

func TestFetchRapidViews(t *testing.T) {

	rapidViews := defaultClient.FetchViews()

	if rapidViews == nil {
		t.Error("rapidViews should not be empty")
	}
}

func _TestSprints(t *testing.T) {

	sprints := defaultClient.FetchSprints(51)

	if sprints == nil {
		t.Error("sprints should not be empty")
	}
}

func _TestSprintDetails(t *testing.T) {

	sprintDetails := defaultClient.FetchSprintDetails(51, 217)

	if sprintDetails == nil {
		t.Error("sprintDetails should not be empty")
	}
}

func _TestSearch(t *testing.T) {

	sprintDetails := defaultClient.FetchSprintDetails(51, 217)

	keys := make([]string, 0, 100)
	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
		keys = append(keys, Issue.Key)
	}
	issues := defaultClient.FetchIssues(keys)

	for _, issue := range issues.Issues {
		for _, item := range issue.Changelog.Histories.Items {
			fmt.Println(issue.Key, issue.Fields.Summary, item.Field, item.FromString, item.ToString)
		}
	}

}
