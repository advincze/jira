package main

import (
	"encoding/json"
	"flag"
	"github.com/advincze/jira-client/jira"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

var port = flag.String("port", "8080", "webserver port")
var openBrowser = flag.Bool("o", false, "open browser at startup")
var testMode = flag.Bool("t", false, "start in test mode")
var cacheMode = flag.Bool("c", false, "cache requests")

func main() {

	flag.Parse()

	if *openBrowser {
		startBrowser("http://localhost:" + *port + "/")
	}
	jira.SetTest(*testMode)
	jira.SetCache(*cacheMode)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/data/burndown", burndownHandler)
	http.HandleFunc("/data/boards", boardsHandler)
	http.HandleFunc("/data/sprints", sprintsHandler)
	http.ListenAndServe(":"+*port, nil)

}

const BOARD_ID = "boardId"
const SPRINT_ID = "sprintId"
const FILTER_LABEL = "filter"

func burndownHandler(w http.ResponseWriter, r *http.Request) {
	boardId, err := strconv.Atoi(r.FormValue(BOARD_ID))
	panicerr(err)

	sprintId, err := strconv.Atoi(r.FormValue(SPRINT_ID))
	panicerr(err)

	sprint := jira.GetSprintById(boardId, sprintId)
	issues := sprint.Issues

	log.Printf("%d issues fetched\n", len(issues))

	if filterLabel := r.FormValue(FILTER_LABEL); filterLabel != "" {
		issues = issues.FilterByLabel(filterLabel)
	}

	log.Printf("%d issues remain after filtering\n", len(issues))
	for _, issue := range issues {
		log.Printf("issue: %s cost: %d\n", issue.Key, issue.EffortInSeconds/3600)
	}

	burndown := jira.CreateBurndown(sprint, issues)

	bytes, _ := json.MarshalIndent(burndown, "", " ")
	w.Write(bytes)
}

func boardsHandler(w http.ResponseWriter, r *http.Request) {
	boards := jira.GetBoards()

	bytes, _ := json.MarshalIndent(boards, "", " ")
	w.Write(bytes)
}

func sprintsHandler(w http.ResponseWriter, r *http.Request) {
	boardId, err := strconv.Atoi(r.FormValue(BOARD_ID))
	panicerr(err)
	board := jira.GetBoardById(boardId)

	bytes, _ := json.MarshalIndent(board.Sprints, "", " ")
	w.Write(bytes)
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
