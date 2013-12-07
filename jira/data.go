package jira

import (
	"errors"
	"fmt"
	"time"
)

type RapidViews struct {
	Views []*RapidView
}

type RapidView struct {
	Id                   int
	Name                 string
	CanEdit              bool
	SprintSupportEnabled bool
}

func (r *RapidViews) GetBoardId(boardName string) (int, error) {
	for _, view := range r.Views {
		if view.Name == boardName {
			return view.Id, nil
		}
	}
	return 0, errors.New("board not found")
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
		}
	}
	return timeline
}
