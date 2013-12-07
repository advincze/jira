package jira

import (
	"encoding/json"
	"testing"
)

func init() {
	// defaultClient.Test = true
}

func TestCreateBurndown(t *testing.T) {

	board := GetBoard("Release Planning Board")

	sprint := board.GetSprint("RC-5.31")

	issues := sprint.Issues //.FilterByLabel("MagicWombats")

	burndown := CreateBurndown(sprint, issues)

	bytes, _ := json.MarshalIndent(&burndown, "", " ")
	t.Logf("%s\n\n", bytes)
}
