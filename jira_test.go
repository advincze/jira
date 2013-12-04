package jira

import (
	"testing"
)

var jiraClient *Jira

func init() {
	config := loadConfig("test.yaml")
	jiraClient = NewJira(config.BaseUrl, config.Login, config.Password)
}

func TestFetchRapidViews(t *testing.T) {

	rapidViews := jiraClient.fetchViews()

	if rapidViews == nil {
		t.Error("rapidViews should not be empty")
	}
}

func TestSprints(t *testing.T) {

	sprints := jiraClient.fetchSprints(51)

	if sprints == nil {
		t.Error("sprints should not be empty")
	}
}

func TestSprintDetails(t *testing.T) {

	sprintDetails := jiraClient.fetchSprintDetails(51, 217)

	if sprintDetails == nil {
		t.Error("sprintDetails should not be empty")
	}
}

func TestSearch(t *testing.T) {

	sprintDetails := jiraClient.fetchSprintDetails(51, 217)

	keys := make([]string, 0, 100)
	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
		keys = append(keys, Issue.Key)
	}

	issues := jiraClient.fetchIssues(keys)

	if issues == nil {
		t.Error("issues should not be empty")
	}
}
