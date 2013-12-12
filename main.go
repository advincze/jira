package main

import (
	"encoding/json"
	"flag"
	"github.com/advincze/jira-client/jira"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
)

var port = flag.String("port", "8080", "webserver port")
var openBrowser = flag.Bool("o", false, "open browser at startup")

func main() {

	flag.Parse()

	if *openBrowser {
		startBrowser("http://localhost:" + *port + "/")
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/data/burndown", burndownHandler)
	http.HandleFunc("/data/boards", boardsHandler)
	http.HandleFunc("/data/sprints", sprintsHandler)
	http.ListenAndServe(":"+*port, nil)

}

func burndownHandler(w http.ResponseWriter, r *http.Request) {
	boardId, err := strconv.Atoi(r.FormValue("board"))
	panicerr(err)

	sprintId, err := strconv.Atoi(r.FormValue("sprint"))
	panicerr(err)

	filter := r.FormValue("team")

	burndown := jira.GetBurndown(boardId, sprintId, filter)

	bytes, _ := json.MarshalIndent(burndown, "", " ")
	w.Write(bytes)
}

func boardsHandler(w http.ResponseWriter, r *http.Request) {
	boards := jira.FetchBoards()
	bytes, _ := json.MarshalIndent(boards, "", " ")
	w.Write(bytes)
}

func sprintsHandler(w http.ResponseWriter, r *http.Request) {
	boardId, err := strconv.Atoi(r.FormValue("board"))
	panicerr(err)

	sprints := jira.FetchSprints(boardId)

	bytes, _ := json.MarshalIndent(sprints, "", " ")
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
