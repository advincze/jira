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

func createBurndownFromSprint(sprintDetails *SprintDetails) *Burndown {
	return createBurndown(sprintDetails.Start, sprintDetails.End, sprintDetails.Issues)
}

func createBurndown(start, end time.Time, issues []*Issue) *Burndown {

	timeline := make([]*TimelineElement, 0, len(issues)*4)

	effortChanges := make(map[time.Time]int)
	totalWorkInHours := 0
	for _, issue := range issues {
		totalWorkInHours += issue.OriginalEstimateSeconds / secondsInHour
		for _, change := range issue.Changes {
			if change.issueClosed() {
				effortChanges[change.Timestamp] = -issue.OriginalEstimateSeconds
				timelineElement := &TimelineElement{
					Timestamp:    change.Timestamp,
					TaskAffected: issue.Key,
				}
				timeline = append(timeline, timelineElement)
			}

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
		SprintStart: start,
		SprintEnd:   end,
		Timeline:    timeline,
	}
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
