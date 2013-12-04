package jira

import (
	"fmt"
	"testing"
)

var jiraClient *Jira

func init() {
	config := loadConfig("test.yaml")
	jiraClient = NewJira(config.BaseUrl, config.Login, config.Password)
	jiraClient.DumpResults = true
}

func TestFetchRapidViews(t *testing.T) {

	rapidViews := jiraClient.FetchViews()

	if rapidViews == nil {
		t.Error("rapidViews should not be empty")
	}
}

func TestSprints(t *testing.T) {

	sprints := jiraClient.FetchSprints(51)

	if sprints == nil {
		t.Error("sprints should not be empty")
	}
}

func TestSprintDetails(t *testing.T) {

	sprintDetails := jiraClient.FetchSprintDetails(51, 217)

	if sprintDetails == nil {
		t.Error("sprintDetails should not be empty")
	}
}

func TestSearch(t *testing.T) {

	sprintDetails := jiraClient.FetchSprintDetails(51, 217)

	keys := make([]string, 0, 100)
	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
		keys = append(keys, Issue.Key)
	}

	issues := jiraClient.fetchIssues(keys)

	// fmt.Printf("%#v\n", issues)

	for _, issue := range issues.Issues {
		for _, item := range issue.Changelog.Histories.Items {
			fmt.Println(issue.Key, issue.Fields.Summary, item.Field, item.FromString, item.ToString)
		}
	}
	t.Error("")
}
