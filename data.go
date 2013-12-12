package jira

type RapidViewsResponse struct {
	Views []*struct {
		Id                   int
		Name                 string
		CanEdit              bool
		SprintSupportEnabled bool
	}
}

func (rv *RapidViewsResponse) getRapidViewNames() []string {
	names := make([]string, 0, len(rv.Views))
	for _, view := range rv.Views {
		names = append(names, view.Name)
	}
	return names
}

func (rv *RapidViewsResponse) getBoardId(boardName string) (int, bool) {
	for _, view := range rv.Views {
		if view.Name == boardName {
			return view.Id, true
		}
	}
	return 0, false
}

type SprintsResponse struct {
	RapidViewId int
	Sprints     []struct {
		Id    int
		Name  string
		State string
	}
}

func (sr *SprintsResponse) getSprintNames() []string {
	names := make([]string, 0, len(sr.Sprints))
	for _, sprint := range sr.Sprints {
		names = append(names, sprint.Name)
	}
	return names
}

func (s *SprintsResponse) getSprintId(sprintName string) (int, bool) {
	for _, sprint := range s.Sprints {
		if sprint.Name == sprintName {
			return sprint.Id, true
		}
	}
	return 0, false
}

type SprintDetailsResponse struct {
	Sprint struct {
		Id        int
		Name      string
		State     string
		StartDate string
		EndDate   string
	}
}

type SearchResult struct {
	Expand string
	Issues []struct {
		Id     string
		Key    string
		Fields struct {
			Created   string
			Issuetype struct {
				Name string
			}
			Labels       []string
			Timetracking struct {
				OriginalEstimateSeconds int
			}
		}
		Changelog struct {
			Histories []History
		}
	}
}

type History struct {
	// Id      int
	Created string
	Items   []struct {
		Field      string
		FromString string
		ToString   string
	}
}

// func (h *History) isClosingEntry() bool {
// 	for _, item := range h.Items {
// 		if item.Field == "status" {
// 			switch item.FromString {
// 			case "Open":
// 				break
// 			case "Planung":
// 				break
// 			case "In Progress":
// 				break
// 			case "Geschlossen":
// 				return false
// 			}

// 			switch item.ToString {
// 			case "Open":
// 				break
// 			case "Planung":
// 				break
// 			case "In Progress":
// 				break
// 			case "Geschlossen":
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// func (s *SearchResult) GetTimeLine(from, to time.Time) []time.Time {
// 	timeline := make([]time.Time, 0, len(s.Issues)*3)
// 	for _, issue := range s.Issues {
// 		if len(issue.Changelog.Histories) == 0 {
// 			fmt.Printf("issue with no changelog: %v, %v, %v \n", issue.Key, issue.Fields.Issuetype.Name, issue.Fields.Timetracking.OriginalEstimateSeconds)
// 		}
// 		for _, history := range issue.Changelog.Histories {
// 			t, _ := time.Parse("2006-01-02T15:04:05.000-0700", history.Created)
// 			if t.After(from) && t.Before(to) {
// 				closed := history.isClosingEntry()
// 				for _, item := range history.Items {
// 					if item.Field == "status" || item.Field == "Sprint" || item.Field == "timeestimate" {
// 						fmt.Printf("%d : %v - [%v] %v -> %v %v , %v \n", t.Unix(), issue.Key, item.Field, item.FromString, item.ToString, closed, issue.Fields.Issuetype.Name)
// 					} else {
// 						fmt.Printf("%d : %v - [%v] \n", t.Unix(), issue.Key, item.Field)
// 					}
// 				}
// 			} else {
// 				for _, item := range history.Items {
// 					fmt.Printf("XXX %d : %v - [%v] \n", t.Unix(), issue.Key, item.Field)
// 				}

// 			}
// 		}
// 	}
// 	return timeline
// }
