package jira

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var cachingTime time.Duration = 3 * time.Minute

type JiraClient struct {
	client *http.Client
	config *Config
}

type Auth struct {
	Login    string
	Password string
}

func NewJira(config *Config) *JiraClient {
	return &JiraClient{
		client: &http.Client{},
		config: config,
	}
}

func NewJiraFromFile(configFile string) *JiraClient {
	config := LoadConfig(configFile)
	return NewJira(config)
}

type Board struct {
	Id   int
	Name string
}

func (jc *JiraClient) FetchBoards() []*Board {
	rapidViewsResponse := jc.fetchRapidViews()
	boards := make([]*Board, 0, len(rapidViewsResponse.Views))
	for _, view := range rapidViewsResponse.Views {
		board := &Board{
			Id:   view.Id,
			Name: view.Name,
		}
		boards = append(boards, board)
	}
	return boards
}

type Sprint struct {
	Id   int
	Name string
}

func (jc *JiraClient) FetchSprints(boardId int) []*Sprint {
	sprintsResponse := jc.fetchSprintResponses(boardId)
	sprints := make([]*Sprint, 0, len(sprintsResponse.Sprints))
	for _, sprint := range sprintsResponse.Sprints {
		sprint := &Sprint{
			Id:   sprint.Id,
			Name: sprint.Name,
		}
		sprints = append(sprints, sprint)
	}
	return sprints
}

type SprintDetails struct {
	Sprint
	Start  time.Time
	End    time.Time
	Issues Issues
}

func (jc *JiraClient) FetchSprintDetails(boardId, sprintId int) *SprintDetails {
	sprintDetailsResponse := jc.fetchSprintDetails(boardId, sprintId)

	start, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetailsResponse.Sprint.StartDate)
	end, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetailsResponse.Sprint.EndDate)

	searchResults := jc.fetchSprintIssues(sprintId)
	issues := make([]*Issue, 0, len(searchResults.Issues))
	for _, foundIssue := range searchResults.Issues {
		changes := make([]*Change, 0, len(foundIssue.Changelog.Histories)*10)
		for _, history := range foundIssue.Changelog.Histories {
			created, _ := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
			for _, item := range history.Items {
				change := &Change{
					Timestamp: created,
					Field:     item.Field,
					From:      item.FromString,
					To:        item.ToString,
				}
				changes = append(changes, change)
			}

		}
		issue := &Issue{
			Id:  foundIssue.Id,
			Key: foundIssue.Key,
			//TODO issue type
			// Type                    :foundIssue.
			Labels:                  foundIssue.Fields.Labels,
			OriginalEstimateSeconds: foundIssue.Fields.Timetracking.OriginalEstimateSeconds,
			Changes:                 changes,
		}

		issues = append(issues, issue)
	}
	return &SprintDetails{
		Sprint: Sprint{
			Id:   sprintId,
			Name: sprintDetailsResponse.Sprint.Name,
		},
		Start:  start,
		End:    end,
		Issues: issues,
	}

}

type Issue struct {
	Id                      string
	Key                     string
	Type                    string
	OriginalEstimateSeconds int
	Labels                  []string
	Changes                 []*Change
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

type Change struct {
	Timestamp time.Time
	Field     string
	From      string
	To        string
}

func (c *Change) issueClosed() bool {

	if c.Field == "status" {
		switch c.From {
		case "Open":
			break
		case "Planung":
			break
		case "In Progress":
			break
		case "Geschlossen":
			return false
		}

		switch c.To {
		case "Open":
			break
		case "Planung":
			break
		case "In Progress":
			break
		case "Geschlossen":

			return true
		}

	}
	return false
}

func (jc *JiraClient) FetchBoardNames() []string {
	rapidViewsResponse := jc.fetchRapidViews()
	return rapidViewsResponse.getRapidViewNames()
}

func (jc *JiraClient) FetchBoardByName(boardName string) *Board {
	rapidViewsResponse := jc.fetchRapidViews()
	if boardId, ok := rapidViewsResponse.getBoardId(boardName); ok {
		return &Board{
			Id:   boardId,
			Name: boardName,
		}
	}
	return nil
}

const (
	rapidViewsEndpoint    = "/rest/greenhopper/1.0/rapidview"
	sprintsEndpoint       = "/rest/greenhopper/1.0/sprintquery/%d"
	sprintDetailsEndpoint = "/rest/greenhopper/1.0/rapid/charts/sprintreport?rapidViewId=%d&sprintId=%d"
	searchEndpoint        = "/rest/api/2/search"
)

func (jc *JiraClient) fetchRapidViews() (rapidViewsResponse *RapidViewsResponse) {
	jc.fetchJson(rapidViewsEndpoint, &rapidViewsResponse)
	return
}

func (jc *JiraClient) fetchSprintResponses(rapidViewId int) (sprintsResponse *SprintsResponse) {
	jc.fetchJson(fmt.Sprintf(sprintsEndpoint, rapidViewId), &sprintsResponse)
	return
}

func (jc *JiraClient) fetchSprintDetails(rapidViewId, sprintId int) (sprintDetailsResponse *SprintDetailsResponse) {
	jc.fetchJson(fmt.Sprintf(sprintDetailsEndpoint, rapidViewId, sprintId), &sprintDetailsResponse)
	return
}

func (jc *JiraClient) fetchSprintIssues(sprintId int) (searchResult *SearchResult) {
	jql := fmt.Sprintf("Sprint=%d", sprintId)
	fields := "*all"
	expand := "changelog"
	//TODO 1000 is a magic number
	return jc.fetchSearchResult(jql, fields, expand, 1000)
}

func (jc *JiraClient) fetchSearchResult(jql, fields, expand string, maxSearchResults int) (searchResult *SearchResult) {
	val := url.Values{}
	val.Set("jql", jql)
	val.Set("fields", fields)
	val.Set("expand", expand)
	val.Set("maxResults", fmt.Sprintf("%d", maxSearchResults))

	jc.fetchJson(searchEndpoint+"?"+val.Encode(), &searchResult)
	return
}

func (jc *JiraClient) fetchJson(endpointUrl string, object interface{}) {

	body := jc.fetchJiraGetRequest(jc.config.BaseUrl + endpointUrl)

	err := json.Unmarshal(body, &object)
	panicerr(err)
}

func gethashFileNameForUrl(url string) string {
	h := sha1.New()
	io.WriteString(h, url)
	return fmt.Sprintf(".cache/%x", h.Sum(nil))
}

var cachedFiles = make(map[string]time.Time)

func (jc *JiraClient) fetchJiraGetRequest(url string) []byte {

	log.Printf("fetching url: [%s]\n", url)

	fileName := gethashFileNameForUrl(url)
	log.Printf("cachefile %s \n", fileName)

	t0 := time.Now().Add(-1 * cachingTime)
	if cachetime, ok := cachedFiles[fileName]; ok && cachetime.After(t0) {
		bytes, err := ioutil.ReadFile(fileName)
		panicerr(err)
		return bytes
	}

	req, err := http.NewRequest("GET", url, nil)
	panicerr(err)
	req.SetBasicAuth(jc.config.Auth.Login, jc.config.Auth.Password)

	resp, err := jc.client.Do(req)
	panicerr(err)

	body, err := ioutil.ReadAll(resp.Body)
	panicerr(err)

	os.Mkdir(".cache", 0777)

	err = ioutil.WriteFile(fileName, body, 0644)
	panicerr(err)

	cachedFiles[fileName] = time.Now()
	log.Printf("adding file %s to cache \n", fileName)
	deleteOutdatedFiles()

	return body
}

func deleteOutdatedFiles() {
	fileinfos, err := ioutil.ReadDir(".cache")
	panicerr(err)
	t := time.Now().Add(-1 * cachingTime)
	for _, fi := range fileinfos {
		filename := ".cache/" + fi.Name()
		if cachetime, ok := cachedFiles[filename]; !ok {
			log.Printf("deleting file not found in cache: %s\n", filename)
			os.Remove(filename)
		} else if cachetime.Before(t) {
			log.Printf("deleting file outdated: %s, cachetime: %v \n", filename, cachetime)
			os.Remove(filename)
			delete(cachedFiles, filename)
		}
	}
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}
