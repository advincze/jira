package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Jira struct {
	BaseUrl     string
	Auth        *Auth
	client      *http.Client
	DumpResults bool
}

type Auth struct {
	login    string
	password string
}

const (
	rapidViewsEndpoint    = "/rest/greenhopper/1.0/rapidview"
	sprintsEndpoint       = "/rest/greenhopper/1.0/sprintquery/%d"
	issueEndpoint         = "/rest/api/2/issue/{issueIdOrKey}"
	searchEndpoint        = "/rest/api/2/search"
	sprintDetailsEndpoint = "/rest/greenhopper/1.0/rapid/charts/sprintreport?rapidViewId=%d&sprintId=%d"
)

func JiraWithConfig(configFile string) *Jira {
	config := LoadConfig(configFile)
	return NewJira(config.BaseUrl, config.Login, config.Password)
}

func NewJira(baseUrl, login, password string) *Jira {
	jira := &Jira{
		BaseUrl: baseUrl,
		Auth: &Auth{
			login:    login,
			password: password,
		},
		client: &http.Client{},
	}
	return jira
}

func (jr *Jira) FetchViews() (rapidViews *RapidViews) {
	jr.fetchJson(rapidViewsEndpoint, &rapidViews)
	return
}

func (jr *Jira) FetchSprints(rapidViewId int) (sprints *Sprints) {
	jr.fetchJson(fmt.Sprintf(sprintsEndpoint, rapidViewId), &sprints)
	return
}

func (jr *Jira) FetchSprintDetails(rapidViewId, sprintId int) (sprintDetails *SprintDetails) {
	jr.fetchJson(fmt.Sprintf(sprintDetailsEndpoint, rapidViewId, sprintId), &sprintDetails)
	return
}

func (jr *Jira) FetchIssues(keys []string) (searchResult *SearchResult) {
	jql := fmt.Sprintf("key in (%s) OR parent in (%[1]s)", strings.Join(keys, ","))
	return jr.SearchIssues(jql, "id", "changelog")
}

func (jr *Jira) SearchIssues(jql, fields, expand string) (searchResult *SearchResult) {
	val := url.Values{}
	val.Set("jql", jql)
	val.Set("fields", fields)
	val.Set("expand", expand)
	jr.fetchJson(searchEndpoint+"?"+val.Encode(), &searchResult)
	return
}

func (jr *Jira) fetchJson(endpointUrl string, object interface{}) {
	req, _ := http.NewRequest("GET", jr.BaseUrl+endpointUrl, nil)

	req.SetBasicAuth(jr.Auth.login, jr.Auth.password)

	resp, _ := jr.client.Do(req)

	if jr.DumpResults {
		file_postfix := time.Now().String()

		bytes, _ := httputil.DumpResponse(resp, false)

		ioutil.WriteFile(file_postfix+".resp", bytes, 0777)

		bytes, _ = httputil.DumpRequest(req, false)

		ioutil.WriteFile(file_postfix+".req", bytes, 0777)

		bytes, _ = ioutil.ReadAll(resp.Body)

		ioutil.WriteFile(file_postfix+".body", bytes, 0777)
	}

	dec := json.NewDecoder(resp.Body)

	dec.Decode(&object)
}

type RapidViews struct {
	Views []*RapidView
}

type RapidView struct {
	Id                   int
	Name                 string
	CanEdit              bool
	SprintSupportEnabled bool
}

type Sprints struct {
	RapidViewId int
	Sprints     []*Sprint
}

type Sprint struct {
	Id        int
	Name      string
	State     string
	StartDate string
	EndDate   string
}

type SprintDetails struct {
	Contents struct {
		CompletedIssues   []*Issue
		IncompletedIssues []*Issue
		PuntedIssues      []*Issue
	}
	Sprint Sprint
}

type Issue struct {
	Id         int
	Key        string
	StatusId   string
	StatusName string
	Expand     string
	Fields     *IssueFields
	Changelog  *Changelog
}

type Changelog struct {
	StartAt   int
	Histories *History
}

type History struct {
	Id    int
	Items []*HistoryItem
}

type HistoryItem struct {
	Field      string
	FromString string
	ToString   string
}

type IssueFields struct {
	Summary     string
	Description string
	Updated     string
	Created     string
	status      struct {
		Name string
	}
	Issuetype *IssueType
	Priority  struct {
		Name string
	}
	Subtasks              []*Issue
	Aggregatetimeestimate int
}

type IssueType struct {
	Self        string
	Id          string
	Description string
	IconUrl     string
	Name        string
	Subtask     bool
}

type SearchResult struct {
	Expand     string
	StartAt    int
	MaxResults int
	Total      int
	Issues     []*Issue
}
