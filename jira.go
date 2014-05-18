package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const maxFetchedIssues = 1000

var cachingTime = 3 * time.Minute

type JiraClient struct {
	client *http.Client
	config *Config
}

type Auth struct {
	Login    string
	Password string
}

func NewJira(config *Config) *JiraClient {
	jira := &JiraClient{
		client: &http.Client{},
		config: config,
	}
	return jira
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
			Id:                      foundIssue.Id,
			Key:                     foundIssue.Key,
			Title:                   foundIssue.Fields.Summary,
			Type:                    foundIssue.Fields.Issuetype.Name,
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
		Issues: Issues(issues),
	}

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
	return jc.fetchSearchResult(jql, fields, expand, maxFetchedIssues)
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

var jiraGetRequestFetcher = func(jc *JiraClient, url string) []byte {

	return defaultCache.getOrRunCacheAndReturn(cachingTime, url, func() interface{} {
		req, err := http.NewRequest("GET", url, nil)
		panicerr(err)
		req.SetBasicAuth(jc.config.Auth.Login, jc.config.Auth.Password)

		resp, err := jc.client.Do(req)
		panicerr(err)

		body, err := ioutil.ReadAll(resp.Body)
		panicerr(err)

		return body
	}).([]byte)
}

func (jc *JiraClient) fetchJiraGetRequest(url string) []byte {
	return jiraGetRequestFetcher(jc, url)
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}
