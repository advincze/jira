package jira

import (
	"fmt"
	"testing"
)

var jiraClient = NewJira("baseurl", "username", "password")

func TestFetchRapidViews(t *testing.T) {

	rapidViews := jiraClient.fetchViews()

	fmt.Printf("%#v \n", rapidViews)
}

func TestSprints(t *testing.T) {

	sprints := jiraClient.fetchSprints(51)

	fmt.Printf("%#v \n", sprints)
}

func TestSprintDetails(t *testing.T) {

	sprintDetails := jiraClient.fetchSprintDetails(51, 217)

	fmt.Printf("%#v \n", sprintDetails)
}

func TestSearch(t *testing.T) {

	keys := make([]string, 0, 100)
	sprintDetails := jiraClient.fetchSprintDetails(51, 217)
	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
		keys = append(keys, Issue.Key)
	}

	issues := jiraClient.fetchIssues(keys)

	fmt.Printf("%#v \n", issues)
}
