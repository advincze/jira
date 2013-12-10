package jira

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type JiraClient struct {
	BaseUrl       string
	Auth          *Auth
	client        *http.Client
	DumpResults   bool
	Test          bool
	CacheRequests bool
}

type Auth struct {
	login    string
	password string
}

var defaultClient *JiraClient

func init() {
	defaultClient = JiraWithConfig("jira.yaml")
}

func SetCache(cache bool) {
	defaultClient.CacheRequests = cache
}

func SetTest(test bool) {
	defaultClient.Test = test
}

func SetConfig(configFile string) {
	defaultClient = JiraWithConfig(configFile)
}

func GetBoard(boardName string) *Board {
	return defaultClient.GetBoard(boardName)
}

func GetBoards() []*Board {
	return defaultClient.GetBoards()
}

func (jiraClient *JiraClient) GetBoards() []*Board {

	views := jiraClient.FetchViews()
	boards := make([]*Board, 0, len(views.Views))
	for _, view := range views.Views {
		boards = append(boards, &Board{
			Id:   view.Id,
			Name: view.Name,
		})
	}

	return boards
}

func (jiraClient *JiraClient) GetBoard(boardName string) *Board {

	views := jiraClient.FetchViews()
	boardId, _ := views.GetBoardId(boardName)
	sprints := jiraClient.FetchSprints(boardId)

	return &Board{
		Id:         boardId,
		Name:       boardName,
		boardId:    boardId,
		boardName:  boardName,
		rapidViews: views,
		Sprints:    sprints,
	}
}

func GetBoardById(boardId int) *Board {
	return defaultClient.GetBoardById(boardId)
}

func (jiraClient *JiraClient) GetBoardById(boardId int) *Board {

	views := jiraClient.FetchViews()
	var boardName string
	for _, view := range views.Views {
		if view.Id == boardId {
			boardName = view.Name
			break
		}
	}

	sprints := jiraClient.FetchSprints(boardId)

	return &Board{
		Id:         boardId,
		Name:       boardName,
		boardId:    boardId,
		boardName:  boardName,
		rapidViews: views,
		Sprints:    sprints,
	}
}

type Board struct {
	Id         int
	Name       string
	boardId    int
	boardName  string
	rapidViews *RapidViews
	Sprints    *Sprints
}

func GetSprintById(boardId, sprintId int) *Sprint {
	log.Printf("getSprintById boardId:%d, sprintId:%d \n", boardId, sprintId)
	sprintDetails := defaultClient.FetchSprintDetails(boardId, sprintId)
	start, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.StartDate)
	end, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.EndDate)

	searchResults := defaultClient.FetchSprintIssues(sprintId)
	// log.Printf("search results %v \n", searchResults)
	// log.Printf("found %d issues to sprintId: %d\n", len(searchResults.Issues), sprintId)
	issues := make([]*Issue, 0, len(searchResults.Issues))
	for _, foundIssue := range searchResults.Issues {
		changes := make([]IssueChange, 0, 10)
		for _, history := range foundIssue.Changelog.Histories {
			created, _ := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
			if history.isClosingEntry() {
				change := IssueChange{
					Timestamp:            created,
					EffortAddedInSeconds: -foundIssue.Fields.Timetracking.OriginalEstimateSeconds,
				}
				changes = append(changes, change)
				// log.Printf("closing issue history: %s  %d \n", foundIssue.Key, len(history.Items))
			} else {
				// log.Printf("non closing issue history: %s  %d \n", foundIssue.Key, len(history.Items))
			}
		}
		issue := &Issue{
			Key:             foundIssue.Key,
			Labels:          foundIssue.Fields.Labels,
			EffortInSeconds: foundIssue.Fields.Timetracking.OriginalEstimateSeconds,
			Changelog:       changes,
		}
		// log.Printf("issue added: %s cost: %d \n", issue.Key, issue.EffortInSeconds/3600)

		issues = append(issues, issue)
	}

	return &Sprint{
		Start:      start,
		End:        end,
		sprintName: sprintDetails.Sprint.Name,
		sprintId:   sprintId,
		Issues:     issues,
	}
}

func (board *Board) GetSprint(sprintName string) *Sprint {
	sprintId, _ := board.Sprints.GetSprintId(sprintName)
	return GetSprintById(board.boardId, sprintId)
}

type Sprint struct {
	Start      time.Time
	End        time.Time
	sprintName string
	sprintId   int
	issueKeys  []string
	Issues     Issues
}

type Issue struct {
	Key             string
	Labels          []string
	EffortInSeconds int
	Changelog       []IssueChange
}

type IssueChange struct {
	Field                string
	Timestamp            time.Time
	EffortAddedInSeconds int
}

type Issues []*Issue

func (issues Issues) FilterByLabel(labelToSearch string) Issues {
	filteredIssues := make([]*Issue, 0, len(issues))
	for _, issue := range issues {
		var containsLabel bool
		for _, labelFound := range issue.Labels {
			if labelFound == labelToSearch {
				containsLabel = true
				break
			}
		}
		if containsLabel {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return Issues(filteredIssues)
}

const (
	rapidViewsEndpoint    = "/rest/greenhopper/1.0/rapidview"
	sprintsEndpoint       = "/rest/greenhopper/1.0/sprintquery/%d"
	issueEndpoint         = "/rest/api/2/issue/{issueIdOrKey}"
	searchEndpoint        = "/rest/api/2/search"
	sprintDetailsEndpoint = "/rest/greenhopper/1.0/rapid/charts/sprintreport?rapidViewId=%d&sprintId=%d"
)

func JiraWithConfig(configFile string) *JiraClient {
	config := LoadConfig(configFile)
	return NewJira(config.BaseUrl, config.Login, config.Password)
}

func NewJira(baseUrl, login, password string) *JiraClient {
	jira := &JiraClient{
		BaseUrl: baseUrl,
		Auth: &Auth{
			login:    login,
			password: password,
		},
		client: &http.Client{},
	}
	return jira
}

func (jr *JiraClient) FetchViews() (rapidViews *RapidViews) {
	if jr.Test {
		jr.loadJsonFile(".testdata/rapidviews.body", &rapidViews)
		return
	}
	jr.fetchJson(rapidViewsEndpoint, &rapidViews)
	return
}

func (jr *JiraClient) FetchSprints(rapidViewId int) (sprints *Sprints) {
	if jr.Test {
		jr.loadJsonFile(".testdata/sprints.body", &sprints)
		return
	}
	jr.fetchJson(fmt.Sprintf(sprintsEndpoint, rapidViewId), &sprints)
	return
}

func (jr *JiraClient) FetchSprintDetails(rapidViewId, sprintId int) (sprintDetails *SprintDetails) {
	if jr.Test {
		jr.loadJsonFile(".testdata/sprintDetails.body", &sprintDetails)
		return
	}
	jr.fetchJson(fmt.Sprintf(sprintDetailsEndpoint, rapidViewId, sprintId), &sprintDetails)
	return
}

func (jr *JiraClient) FetchIssues(keys []string) (searchResult *SearchResult) {
	jql := fmt.Sprintf("key in (%s) OR parent in (%[1]s)", strings.Join(keys, ","))
	return jr.SearchIssues(jql, "*all", "changelog")
}

func (jr *JiraClient) FetchSprintIssues(sprintId int) (searchResult *SearchResult) {
	if jr.Test {
		jr.loadJsonFile(".testdata/sprintIssues.body", &searchResult)
		return
	}
	jql := fmt.Sprintf("Sprint=%d", sprintId)
	return jr.SearchIssues(jql, "*all", "changelog")
}

const maxSearchResults = 1000

func (jr *JiraClient) SearchIssues(jql, fields, expand string) (searchResult *SearchResult) {
	val := url.Values{}
	val.Set("jql", jql)
	val.Set("fields", fields)
	val.Set("expand", expand)
	val.Set("maxResults", fmt.Sprintf("%d", maxSearchResults))

	jr.fetchJson(searchEndpoint+"?"+val.Encode(), &searchResult)
	return
}

func (jr *JiraClient) loadJsonFile(fileName string, object interface{}) {

	bytes, _ := ioutil.ReadFile(fileName)

	json.Unmarshal(bytes, &object)
}

func (jr *JiraClient) fetchJson(endpointUrl string, object interface{}) {

	var fileName string
	if jr.CacheRequests {
		log.Printf("caching on... ")
		h := sha1.New()
		io.WriteString(h, endpointUrl)

		fileName = fmt.Sprintf(".cache/%x", h.Sum(nil))
		log.Printf("looking if file exists:  %s in cache", fileName)

		if _, err := os.Stat(fileName); err == nil {

			log.Printf("trying to read file %s from cache", fileName)
			if bytes, err := ioutil.ReadFile(fileName); err == nil {
				json.Unmarshal(bytes, &object)
				return
			} else {
				log.Printf("cannot read file %s : %v \n", err)
			}
		}

	}

	url := jr.BaseUrl + endpointUrl
	log.Printf("fetchJson: %s\n", url)

	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(jr.Auth.login, jr.Auth.password)

	resp, _ := jr.client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	if jr.DumpResults {
		file_postfix := time.Now().String()

		bytes, _ := httputil.DumpResponse(resp, false)

		ioutil.WriteFile(file_postfix+".resp", bytes, 0644)

		bytes, _ = httputil.DumpRequest(req, false)

		ioutil.WriteFile(file_postfix+".req", bytes, 0644)

		ioutil.WriteFile(file_postfix+".body", body, 0644)
	}

	if jr.CacheRequests {

		// fileName := ".cache/" + base32.StdEncoding.EncodeToString([]byte(endpointUrl))
		log.Printf("fetched : %d bytes\n", len(body))

		os.Mkdir(".cache", 0777)

		err := ioutil.WriteFile(fileName, body, 0644)
		if err != nil {
			log.Println(err)
		}
	}

	json.Unmarshal(body, &object)

}
