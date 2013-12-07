package main

import (
	"fmt"
	"github.com/advincze/jira-client/jira"
	"time"
)

func main2() {

	// boardName := "Release Planning Board"
	// sprintName := "RC-5.33"
	// // label := "MagicWobmbats"

	jiraClient := jira.JiraWithConfig("test.yaml")

	// views := jiraClient.FetchViews()
	// boardId, _ := views.GetBoardId(boardName)

	// fmt.Printf("boardid: %d \n", boardId)

	// sprints := jiraClient.FetchSprints(boardId)
	// sprintId, _ := sprints.GetSprintId(sprintName)

	// fmt.Printf("sprintid: %d \n", sprintId)

	// sprintDetails := jiraClient.FetchSprintDetails(boardId, sprintId)
	// jiraClient.DumpResults = true
	sprintDetails := jiraClient.FetchSprintDetails(51, 217)
	keys := sprintDetails.GetIssueKeys()

	// jiraClient.DumpResults = true
	issues := jiraClient.FetchIssues(keys)
	println(issues.Total)
	start, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.StartDate)
	end, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.EndDate)
	// Mon Jan 2 15:04:05 -0700 MST 2006
	fmt.Printf("start: %v, end: %v\n", start, end)
	issues.GetTimeLine(start, end)

	// fmt.Printf("timeline: %v\n", issues.GetTimeLine())
	//

	// issues := jiraClient.FetchIssues(keys)

	// // fmt.Printf("%#v\n", issues)

	// for _, issue := range issues.Issues {
	// 	fmt.Println(issue.Key, issue.Changelog)
	// 	// for _, item := range issue.Changelog.Histories.Items {
	// 	// 	fmt.Println(issue.Key, issue.Fields.Summary, item.Field, item.FromString, item.ToString)
	// 	// }
	// }

}

type Burndown struct {
	SprintStart  time.Time
	SprintEnd    time.Time
	IdealLine    map[time.Time]int
	RealLine     map[time.Time]int
	TaskAffected map[time.Time]string
}

func main() {

	board := jira.GetBoard("Release Planning Board")
	fmt.Printf("%v\n", board)
	sprint := board.GetSprint("RC-5.33")
	fmt.Printf("%v\n", sprint)
	issues := sprint.GetIssues().FilterByLabel("MagicWombats")
	fmt.Printf("%v\n", issues)
	burndown := sprint.CreateBurndown(issues)
	// burndown2 := sprint.CreateDailyBurndown(issues)
}
