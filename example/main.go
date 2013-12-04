package main

import (
	"fmt"
	"github.com/advincze/jira-client"
)

func main() {

	config := jira.LoadConfig("test.yaml")
	jiraClient := jira.NewJira(config.BaseUrl, config.Login, config.Password)
	jiraClient.DumpResults = true

	sprintDetails := jiraClient.FetchSprintDetails(51, 217)

	keys := make([]string, 0, 100)
	for _, Issue := range sprintDetails.Contents.IncompletedIssues {
		keys = append(keys, Issue.Key)
	}

	issues := jiraClient.FetchIssues(keys)

	// fmt.Printf("%#v\n", issues)

	for _, issue := range issues.Issues {
		fmt.Println(issue.Key, issue.Changelog)
		// for _, item := range issue.Changelog.Histories.Items {
		// 	fmt.Println(issue.Key, issue.Fields.Summary, item.Field, item.FromString, item.ToString)
		// }
	}

}
