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
	BaseUrl string
	Auth    *Auth
	client  *http.Client
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

func (jr *Jira) fetchViews() (rapidViews RapidViews) {
	jr.fetchJson(rapidViewsEndpoint, &rapidViews)
	return
}

func (jr *Jira) fetchSprints(rapidViewId int) (sprints Sprints) {
	jr.fetchJson(fmt.Sprintf(sprintsEndpoint, rapidViewId), &sprints)
	return
}

func (jr *Jira) fetchSprintDetails(rapidViewId, sprintId int) (sprintDetails SprintDetails) {
	jr.fetchJson(fmt.Sprintf(sprintDetailsEndpoint, rapidViewId, sprintId), &sprintDetails)
	return
}

func (jr *Jira) fetchIssues(keys []string) (searchResult interface{}) {
	jql := fmt.Sprintf("key in (%s)", strings.Join(keys, ","))
	println(jql)
	return jr.SearchIssues(jql, "*navigable", "")
}

func (jr *Jira) SearchIssues(jql, fields, expand string) (searchResult interface{}) {
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

	bytes, _ := httputil.DumpResponse(resp, true)
	fmt.Println("bytes", string(bytes))

	ioutil.WriteFile("file"+time.Now().String()+".out", bytes, 0777)

	dec := json.NewDecoder(resp.Body)

	dec.Decode(&object)
}

type RapidViews struct {
	Views []RapidView
}

type RapidView struct {
	Id                   int
	Name                 string
	CanEdit              bool
	SprintSupportEnabled bool
}

type Sprints struct {
	RapidViewId int
	Sprints     []Sprint
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
		CompletedIssues   []Issue
		IncompletedIssues []Issue
		PuntedIssues      []Issue
	}
	Sprint Sprint
}

type Issue struct {
	Id         int
	Key        string
	StatusId   string
	StatusName string
	fields     struct {
	}
}

type SearchResult struct {
	Issues []Issue
}
