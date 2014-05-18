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
			Summary   string
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
