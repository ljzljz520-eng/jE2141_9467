package report

import (
	"fmt"
	"sort"

	"volunteerhours/domain"
	"volunteerhours/query"
)

type Summary struct {
	RecordCount    int
	VolunteerCount int
	ActivityCount  int
	TotalHours     float64
	PendingCount   int
	ApprovedCount  int
	RejectedCount  int
	ByVolunteer    []query.Group
	ByActivity     []query.Group
}

func BuildSummary(result query.Result, volunteers []domain.Volunteer, activities []domain.Activity) Summary {
	summary := Summary{RecordCount: result.Count(), VolunteerCount: len(volunteers), ActivityCount: len(activities), TotalHours: result.Hours()}
	for _, record := range result.Records {
		switch record.Status {
		case domain.RecordPending:
			summary.PendingCount++
		case domain.RecordApproved:
			summary.ApprovedCount++
		case domain.RecordRejected:
			summary.RejectedCount++
		}
	}
	summary.ByVolunteer = query.GroupByVolunteer(result)
	summary.ByActivity = query.GroupByActivity(result)
	return summary
}

func (s Summary) StatusLabel() string {
	return fmt.Sprintf("pending=%d approved=%d rejected=%d", s.PendingCount, s.ApprovedCount, s.RejectedCount)
}

func (s Summary) TopVolunteer() (query.Group, bool) {
	if len(s.ByVolunteer) == 0 {
		return query.Group{}, false
	}
	groups := append([]query.Group(nil), s.ByVolunteer...)
	sort.SliceStable(groups, func(i, j int) bool { return groups[i].TotalHour > groups[j].TotalHour })
	return groups[0], true
}
