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

type JiraClient struct {
	BaseUrl     string
	Auth        *Auth
	client      *http.Client
	DumpResults bool
	Test        bool
}

type Auth struct {
	login    string
	password string
}

var defaultClient *JiraClient

func init() {
	defaultClient = JiraWithConfig("jira.yaml")
}

func SetConfig(configFile string) {
	defaultClient = JiraWithConfig(configFile)
}

func GetBoard(boardName string) *Board {
	return defaultClient.GetBoard(boardName)
}

func (jiraClient *JiraClient) GetBoard(boardName string) *Board {

	views := jiraClient.FetchViews()
	boardId, _ := views.GetBoardId(boardName)
	sprints := jiraClient.FetchSprints(boardId)

	return &Board{
		boardId:    boardId,
		boardName:  boardName,
		rapidViews: views,
		sprints:    sprints,
	}
}

type Board struct {
	boardId    int
	boardName  string
	rapidViews *RapidViews
	sprints    *Sprints
}

func (b *Board) GetSprint(sprintName string) *Sprint {
	sprintId, _ := b.sprints.GetSprintId(sprintName)
	sprintDetails := defaultClient.FetchSprintDetails(b.boardId, sprintId)
	// keys := sprintDetails.GetIssueKeys()
	issuesx := defaultClient.FetchSprintIssues(sprintId)
	start, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.StartDate)
	end, _ := time.Parse("02/Jan/06 15:04 PM", sprintDetails.Sprint.EndDate)
	issues := make([]*Issue, 0, len(issuesx.Issues))
	for _, issuex := range issuesx.Issues {

		issue := &Issue{
			Key:    issuex.Key,
			Labels: issuex.Fields.Labels,
		}
		issues = append(issues, issue)
	}

	return &Sprint{
		Start:      start,
		End:        end,
		sprintName: sprintName,
		sprintId:   sprintId,
		Issues:     issues,
	}
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
	Key    string
	Labels []string
}

type Issues []*Issue

func (issues Issues) FilterByLabel(labelToSearch string) Issues {
	filteredIssues := make([]*Issue, 0, len(issues))
	// fmt.Printf("filtering: %d issues \n", len(issues))
	for _, issue := range issues {
		// fmt.Printf("issue: %v \n", issue.Key)
		var containsLabel bool
		for _, labelFound := range issue.Labels {
			// fmt.Printf("label found: %s \n", labelFound)
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

func (jr *JiraClient) SearchIssues(jql, fields, expand string) (searchResult *SearchResult) {
	val := url.Values{}
	val.Set("jql", jql)
	val.Set("fields", fields)
	val.Set("expand", expand)

	jr.fetchJson(searchEndpoint+"?"+val.Encode(), &searchResult)
	return
}

func (jr *JiraClient) loadJsonFile(fileName string, object interface{}) {

	bytes, _ := ioutil.ReadFile(fileName)

	json.Unmarshal(bytes, &object)
}

func (jr *JiraClient) fetchJson(endpointUrl string, object interface{}) {
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
