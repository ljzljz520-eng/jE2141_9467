package workflow

import (
	"sort"

	"volunteerhours/domain"
	"volunteerhours/query"
)

type Metrics struct {
	TotalRecords       int
	DistinctVolunteers int
	DistinctActivities int
	TotalHours         float64
	ApprovedHours      float64
	PendingHours       float64
	RejectedHours      float64
	Daily              []query.DailyTotal
	Statuses           []query.StatusTotal
}

func BuildMetrics(records []domain.ServiceRecord) Metrics {
	return Metrics{
		TotalRecords:       len(records),
		DistinctVolunteers: len(query.DistinctVolunteers(records)),
		DistinctActivities: len(query.DistinctActivities(records)),
		TotalHours:         query.Result{Records: records}.Hours(),
		ApprovedHours:      query.ApprovedHours(records),
		PendingHours:       query.PendingHours(records),
		RejectedHours:      query.RejectedHours(records),
		Daily:              query.DailyTotals(records),
		Statuses:           query.StatusTotals(records),
	}
}

func (m Metrics) IsEmpty() bool {
	return m.TotalRecords == 0 && m.TotalHours == 0
}

func (m Metrics) ApprovalRate() float64 {
	if m.TotalRecords == 0 {
		return 0
	}
	approved := 0
	for _, status := range m.Statuses {
		if status.Status == domain.RecordApproved {
			approved = status.Count
		}
	}
	return float64(approved) / float64(m.TotalRecords)
}

func (m Metrics) PeakDay() (query.DailyTotal, bool) {
	if len(m.Daily) == 0 {
		return query.DailyTotal{}, false
	}
	result := append([]query.DailyTotal(nil), m.Daily...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Hours == result[j].Hours {
			return result[i].Date < result[j].Date
		}
		return result[i].Hours > result[j].Hours
	})
	return result[0], true
}

func (m Metrics) StatusCount(status string) int {
	for _, item := range m.Statuses {
		if item.Status == status {
			return item.Count
		}
	}
	return 0
}

func (m Metrics) HoursForStatus(status string) float64 {
	for _, item := range m.Statuses {
		if item.Status == status {
			return item.Hours
		}
	}
	return 0
}
