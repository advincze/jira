package jira

import (
	"time"
)

type Burndown struct {
	SprintStart  time.Time
	SprintEnd    time.Time
	IdealLine    map[time.Time]int
	RealLine     map[time.Time]int
	TaskAffected map[time.Time]string
}

func CreateBurndown(sprint *Sprint, issues []*Issue) *Burndown {

	// for _, issue := range issues {

	// }

	return &Burndown{
		SprintStart: sprint.Start,
		SprintEnd:   sprint.End,
	}
}
