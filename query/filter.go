package query

import (
	"strings"

	"volunteerhours/domain"
)

type Filter struct {
	VolunteerID string
	ActivityID  string
	Status      string
	StartDate   string
	EndDate     string
	MinHours    float64
	MaxHours    float64
	RecordedBy  string
	Text        string
}

type Result struct {
	Records    []domain.ServiceRecord
	Volunteers map[string]domain.Volunteer
	Activities map[string]domain.Activity
}

func (f Filter) IsEmpty() bool {
	return strings.TrimSpace(f.VolunteerID) == "" && strings.TrimSpace(f.ActivityID) == "" && strings.TrimSpace(f.Status) == "" && strings.TrimSpace(f.StartDate) == "" && strings.TrimSpace(f.EndDate) == "" && f.MinHours == 0 && f.MaxHours == 0 && strings.TrimSpace(f.RecordedBy) == "" && strings.TrimSpace(f.Text) == ""
}

func ApplyFilter(records []domain.ServiceRecord, filter Filter) []domain.ServiceRecord {
	result := make([]domain.ServiceRecord, 0, len(records))
	for _, record := range records {
		if !matches(record, filter) {
			continue
		}
		result = append(result, record)
	}
	return result
}

func matches(record domain.ServiceRecord, filter Filter) bool {
	if filter.VolunteerID != "" && record.VolunteerID != filter.VolunteerID {
		return false
	}
	if filter.ActivityID != "" && record.ActivityID != filter.ActivityID {
		return false
	}
	if filter.Status != "" && record.Status != filter.Status {
		return false
	}
	if filter.StartDate != "" && record.ServiceDate < filter.StartDate {
		return false
	}
	if filter.EndDate != "" && record.ServiceDate > filter.EndDate {
		return false
	}
	if filter.MinHours > 0 && record.Hours < filter.MinHours {
		return false
	}
	if filter.MaxHours > 0 && record.Hours > filter.MaxHours {
		return false
	}
	if filter.RecordedBy != "" && !strings.EqualFold(record.RecordedBy, filter.RecordedBy) {
		return false
	}
	if filter.Text != "" && !strings.Contains(strings.ToLower(record.Notes), strings.ToLower(filter.Text)) {
		return false
	}
	return true
}

func (r Result) Count() int {
	return len(r.Records)
}

func (r Result) Hours() float64 {
	total := 0.0
	for _, record := range r.Records {
		if record.Status != domain.RecordRejected {
			total += record.Hours
		}
	}
	return total
}

func (r Result) VolunteerName(id string) string {
	if volunteer, ok := r.Volunteers[id]; ok {
		return volunteer.DisplayName()
	}
	return id
}

func (r Result) ActivityTitle(id string) string {
	if activity, ok := r.Activities[id]; ok {
		return activity.Title
	}
	return id
}
