package main

import (
	"fmt"
	"github.com/advincze/jira-client/jira"
)

func main() {

	board := jira.GetBoard("Release Planning Board")
	// fmt.Printf("%q\n\n", board)
	sprint := board.GetSprint("RC-5.33")
	// fmt.Printf("%q\n\n", sprint)
	issues := sprint.GetIssues().FilterByLabel("MagicWombats")
	// fmt.Printf("%q\n\n", issues)
	burndown := jira.CreateBurndown(sprint, issues)
	fmt.Printf("%q\n\n", burndown)
}
