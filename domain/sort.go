package domain

import (
	"sort"
	"strings"
)

type SortField string

const (
	SortByDate      SortField = "date"
	SortByHours     SortField = "hours"
	SortByVolunteer SortField = "volunteer"
	SortByActivity  SortField = "activity"
	SortByStatus    SortField = "status"
)

type SortSpec struct {
	Field      SortField
	Descending bool
}

func ParseSortSpec(field string, descending bool) SortSpec {
	switch strings.ToLower(strings.TrimSpace(field)) {
	case string(SortByHours):
		return SortSpec{Field: SortByHours, Descending: descending}
	case string(SortByVolunteer):
		return SortSpec{Field: SortByVolunteer, Descending: descending}
	case string(SortByActivity):
		return SortSpec{Field: SortByActivity, Descending: descending}
	case string(SortByStatus):
		return SortSpec{Field: SortByStatus, Descending: descending}
	default:
		return SortSpec{Field: SortByDate, Descending: descending}
	}
}

func SortRecords(records []ServiceRecord, spec SortSpec, volunteers map[string]Volunteer, activities map[string]Activity) []ServiceRecord {
	result := append([]ServiceRecord(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		left, right := result[i], result[j]
		var less bool
		switch spec.Field {
		case SortByHours:
			less = left.Hours < right.Hours
		case SortByVolunteer:
			less = strings.ToLower(volunteers[left.VolunteerID].DisplayName()) < strings.ToLower(volunteers[right.VolunteerID].DisplayName())
		case SortByActivity:
			less = strings.ToLower(activities[left.ActivityID].Title) < strings.ToLower(activities[right.ActivityID].Title)
		case SortByStatus:
			less = left.Status < right.Status
		default:
			less = left.ServiceDate < right.ServiceDate
		}
		if left == right {
			return false
		}
		if spec.Descending {
			return !less && left.ID > right.ID
		}
		return less || (left.ServiceDate == right.ServiceDate && left.ID < right.ID)
	})
	return result
}
