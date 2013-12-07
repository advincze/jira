package jira

import (
	"encoding/json"
	"errors"
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

var jiraClient = JiraWithConfig("test.yaml")

func GetBoard(boardName string) *Board {

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
	sprintDetails := jiraClient.FetchSprintDetails(b.boardId, sprintId)
	keys := sprintDetails.GetIssueKeys()
	issuesx := jiraClient.FetchIssues(keys)
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
		sprintx:    sprintDetails.Sprint,
		Start:      start,
		End:        end,
		sprintName: sprintName,
		sprintId:   sprintId,
		issueKeys:  keys,
		issuesx:    issuesx.Issues,
		Issues:     issues,
	}
}

type Sprint struct {
	sprintx    *SprintX
	Start      time.Time
	End        time.Time
	sprintName string
	sprintId   int
	issueKeys  []string
	issuesx    []*IssueX
	Issues     []*Issue
}

func (s *Sprint) GetIssues() Issues {
	return Issues(s.Issues)
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
	return jr.SearchIssues(jql, "*all", "changelog")
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

func (r *RapidViews) GetBoardId(boardName string) (int, error) {
	for _, view := range r.Views {
		if view.Name == boardName {
			return view.Id, nil
		}
	}
	return 0, errors.New("board not found")
}

type RapidView struct {
	Id                   int
	Name                 string
	CanEdit              bool
	SprintSupportEnabled bool
}

type Sprints struct {
	RapidViewId int
	Sprints     []*SprintX
}

func (s *Sprints) GetSprintId(sprintName string) (int, error) {
	for _, sprint := range s.Sprints {
		if sprint.Name == sprintName {
			return sprint.Id, nil
		}
	}
	return 0, errors.New("sprint not found")
}

type SprintX struct {
	Id        int
	Name      string
	State     string
	StartDate string
	EndDate   string
}

type SprintDetails struct {
	Contents struct {
		CompletedIssues   []*IssueX
		IncompletedIssues []*IssueX
		PuntedIssues      []*IssueX
	}
	Sprint *SprintX
}

func (s *SprintDetails) GetIssueKeys() []string {
	keys := make([]string, 0, len(s.Contents.CompletedIssues)+len(s.Contents.IncompletedIssues))
	for _, issue := range s.Contents.CompletedIssues {
		keys = append(keys, issue.Key)
	}
	for _, issue := range s.Contents.IncompletedIssues {
		keys = append(keys, issue.Key)
	}
	for _, issue := range s.Contents.PuntedIssues {
		keys = append(keys, issue.Key)
	}
	return keys
}

type IssueX struct {
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
	Histories []*History
}

type History struct {
	Id      int
	Created string
	Items   []*HistoryItem
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
	Subtasks              []*IssueX
	Aggregatetimeestimate int
	Labels                []string
	Timetracking          struct {
		OriginalEstimateSeconds int
	}
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
	Issues     []*IssueX
}

func Closed(history *History) bool {
	for _, item := range history.Items {
		if item.Field == "status" {
			switch item.FromString {
			case "Open":
				break
			case "Planung":
				break
			case "In Progress":
				break
			case "Geschlossen":
				return false
			}

			switch item.ToString {
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
	}
	return false

}

func (s *SearchResult) GetTimeLine(from, to time.Time) []time.Time {
	timeline := make([]time.Time, 0, len(s.Issues)*3)
	for _, issue := range s.Issues {
		if len(issue.Changelog.Histories) == 0 {
			fmt.Printf("issue with no changelog: %v, %v, %v \n", issue.Key, issue.Fields.Issuetype.Name, issue.Fields.Timetracking.OriginalEstimateSeconds)
		}
		for _, history := range issue.Changelog.Histories {
			t, _ := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
			if t.After(from) && t.Before(to) {
				closed := Closed(history)
				for _, item := range history.Items {
					if item.Field == "status" || item.Field == "Sprint" || item.Field == "timeestimate" {
						fmt.Printf("%d : %v - [%v] %v -> %v %v , %v \n", t.Unix(), issue.Key, item.Field, item.FromString, item.ToString, closed, issue.Fields.Issuetype.Name)
					} else {
						fmt.Printf("%d : %v - [%v] \n", t.Unix(), issue.Key, item.Field)
					}
				}
			} else {
				for _, item := range history.Items {
					fmt.Printf("XXX %d : %v - [%v] \n", t.Unix(), issue.Key, item.Field)
				}

			}

			// if Closed(history) {
			// 	t, _ := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
			// 	fmt.Printf("at time %v the estimate of %v changed by %v \n", t.String(), issue.Key, issue.Fields.Timetracking.OriginalEstimateSeconds)
			// }

		}
	}
	return timeline
}
