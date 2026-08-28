package query

import (
	"sort"

	"volunteerhours/domain"
)

type DailyTotal struct {
	Date  string
	Hours float64
	Count int
}

type StatusTotal struct {
	Status string
	Hours  float64
	Count  int
}

func DailyTotals(records []domain.ServiceRecord) []DailyTotal {
	byDate := make(map[string]*DailyTotal)
	for _, record := range records {
		entry, ok := byDate[record.ServiceDate]
		if !ok {
			entry = &DailyTotal{Date: record.ServiceDate}
			byDate[record.ServiceDate] = entry
		}
		entry.Count++
		if record.Status != domain.RecordRejected {
			entry.Hours += record.Hours
		}
	}
	result := make([]DailyTotal, 0, len(byDate))
	for _, entry := range byDate {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result
}

func StatusTotals(records []domain.ServiceRecord) []StatusTotal {
	byStatus := make(map[string]*StatusTotal)
	for _, record := range records {
		entry, ok := byStatus[record.Status]
		if !ok {
			entry = &StatusTotal{Status: record.Status}
			byStatus[record.Status] = entry
		}
		entry.Count++
		entry.Hours += record.Hours
	}
	result := make([]StatusTotal, 0, len(byStatus))
	for _, entry := range byStatus {
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Status < result[j].Status })
	return result
}

func ApprovedHours(records []domain.ServiceRecord) float64 {
	total := 0.0
	for _, record := range records {
		if record.Status == domain.RecordApproved {
			total += record.Hours
		}
	}
	return total
}

func PendingHours(records []domain.ServiceRecord) float64 {
	total := 0.0
	for _, record := range records {
		if record.Status == domain.RecordPending {
			total += record.Hours
		}
	}
	return total
}

func RejectedHours(records []domain.ServiceRecord) float64 {
	total := 0.0
	for _, record := range records {
		if record.Status == domain.RecordRejected {
			total += record.Hours
		}
	}
	return total
}

func DistinctVolunteers(records []domain.ServiceRecord) []string {
	seen := make(map[string]bool)
	for _, record := range records {
		seen[record.VolunteerID] = true
	}
	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

func DistinctActivities(records []domain.ServiceRecord) []string {
	seen := make(map[string]bool)
	for _, record := range records {
		seen[record.ActivityID] = true
	}
	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}
