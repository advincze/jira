package jira

import (
	"sort"
	"time"
)

type Burndown struct {
	SprintStart time.Time          `json:"sprintstart"`
	SprintEnd   time.Time          `json:"sprintend"`
	Timeline    []*TimelineElement `json:"timeline"`
}

type TimelineElement struct {
	Timestamp            time.Time `json:"timestamp"`
	TotalWorkInHours     int       `json:"totalWorkInHours"`
	RemainingWorkInHours int       `json:"remainingWorkInHours"`
	TaskAffected         string    `json:"taskAffected"`
}

const secondsInHour = 3600

func CreateBurndown(sprint *Sprint, issues []*Issue) *Burndown {

	timeline := make([]*TimelineElement, 0, len(issues)*4)

	effortChanges := make(map[time.Time]int)
	totalWorkInHours := 0
	for _, issue := range issues {
		totalWorkInHours += issue.EffortInSeconds / secondsInHour
		for _, change := range issue.Changelog {
			effortChanges[change.Timestamp] = change.EffortAddedInSeconds
			timelineElement := &TimelineElement{
				Timestamp:            change.Timestamp,
				TotalWorkInHours:     0,
				RemainingWorkInHours: 0,
				TaskAffected:         issue.Key,
			}
			timeline = append(timeline, timelineElement)
		}
	}

	sort.Sort(ByTimestamp{timeline})

	sumSeconds := 0
	for _, element := range timeline {
		t := element.Timestamp
		sumSeconds += effortChanges[t]
		element.TotalWorkInHours = totalWorkInHours
		element.RemainingWorkInHours = totalWorkInHours + sumSeconds/secondsInHour
	}

	return &Burndown{
		SprintStart: sprint.Start,
		SprintEnd:   sprint.End,
		Timeline:    timeline,
	}
}

type Times []time.Time

func (tt Times) Len() int {
	return len(tt)
}

func (tt Times) Swap(i, j int) {
	tt[i], tt[j] = tt[j], tt[i]
}

func (tt Times) Less(i, j int) bool {
	return tt[i].Before(tt[j])
}

type TimelineElements []*TimelineElement

func (tt TimelineElements) Len() int {
	return len(tt)
}

func (tt TimelineElements) Swap(i, j int) {
	tt[i], tt[j] = tt[j], tt[i]
}

type ByTimestamp struct{ TimelineElements }

func (tt ByTimestamp) Less(i, j int) bool {
	return tt.TimelineElements[i].Timestamp.Before(tt.TimelineElements[j].Timestamp)
}
