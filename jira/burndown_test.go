package jira

import (
	"encoding/json"
	"testing"
)

func init() {
	testClient = NewJiraFromFile("../jira.yaml")
}

func TestCreateBurndown(t *testing.T) {

	sprint := FetchSprint("Release Planning Board", "RC-5.34")

	issues := Issues(sprint.Issues).FilterByLabel("MagicWombats")

	burndown := CreateBurndown(sprint, issues)

	bytes, _ := json.MarshalIndent(&burndown, "", " ")
	t.Logf("%s\n\n", bytes)
}
