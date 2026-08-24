package service

import (
	"fmt"
	"sort"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/query"
)

type AnalyticsService struct {
	catalog *CatalogService
	records *RecordService
	queries *QueryService
}

type VolunteerScore struct {
	VolunteerID string
	Name        string
	Hours       float64
	Records     int
	Approved    int
}

type ActivityScore struct {
	ActivityID string
	Title      string
	Hours      float64
	Volunteers int
	Capacity   int
}

func NewAnalyticsService(catalog *CatalogService, records *RecordService, queries *QueryService) *AnalyticsService {
	return &AnalyticsService{catalog: catalog, records: records, queries: queries}
}

func (s *AnalyticsService) VolunteerScores() ([]VolunteerScore, error) {
	volunteers, err := s.catalog.ListVolunteers()
	if err != nil {
		return nil, err
	}
	records, err := s.records.ListRecords()
	if err != nil {
		return nil, err
	}
	scores := make(map[string]*VolunteerScore, len(volunteers))
	for _, volunteer := range volunteers {
		scores[volunteer.ID] = &VolunteerScore{VolunteerID: volunteer.ID, Name: volunteer.DisplayName()}
	}
	for _, record := range records {
		score, ok := scores[record.VolunteerID]
		if !ok {
			continue
		}
		score.Records++
		if record.Status != domain.RecordRejected {
			score.Hours += record.Hours
		}
		if record.Status == domain.RecordApproved {
			score.Approved++
		}
	}
	result := make([]VolunteerScore, 0, len(scores))
	for _, score := range scores {
		result = append(result, *score)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Hours == result[j].Hours {
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		}
		return result[i].Hours > result[j].Hours
	})
	return result, nil
}

func (s *AnalyticsService) ActivityScores() ([]ActivityScore, error) {
	activities, err := s.catalog.ListActivities()
	if err != nil {
		return nil, err
	}
	records, err := s.records.ListRecords()
	if err != nil {
		return nil, err
	}
	scores := make(map[string]*ActivityScore, len(activities))
	volunteers := make(map[string]map[string]bool)
	for _, activity := range activities {
		scores[activity.ID] = &ActivityScore{ActivityID: activity.ID, Title: activity.Title, Capacity: activity.Capacity}
		volunteers[activity.ID] = make(map[string]bool)
	}
	for _, record := range records {
		score, ok := scores[record.ActivityID]
		if !ok {
			continue
		}
		if record.Status != domain.RecordRejected {
			score.Hours += record.Hours
		}
		volunteers[record.ActivityID][record.VolunteerID] = true
	}
	result := make([]ActivityScore, 0, len(scores))
	for id, score := range scores {
		score.Volunteers = len(volunteers[id])
		result = append(result, *score)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Title < result[j].Title })
	return result, nil
}

func (s *AnalyticsService) SearchText(text string) (query.Result, error) {
	return s.queries.Search(query.Filter{Text: strings.TrimSpace(text)}, domain.ParseSortSpec("date", false))
}

func (s *AnalyticsService) ValidateCapacity(activityID string) (int, int, error) {
	activity, err := s.catalog.FindActivity(activityID)
	if err != nil {
		return 0, 0, err
	}
	records, err := s.records.ListRecords()
	if err != nil {
		return 0, 0, err
	}
	people := make(map[string]bool)
	for _, record := range records {
		if record.ActivityID == activityID && record.Status != domain.RecordRejected {
			people[record.VolunteerID] = true
		}
	}
	return len(people), activity.Capacity, nil
}

func (s *AnalyticsService) CapacityMessage(activityID string) (string, error) {
	used, capacity, err := s.ValidateCapacity(activityID)
	if err != nil {
		return "", err
	}
	if capacity == 0 {
		return "no capacity configured", nil
	}
	if used >= capacity {
		return fmt.Sprintf("capacity reached (%d/%d)", used, capacity), nil
	}
	return fmt.Sprintf("capacity available (%d/%d)", used, capacity), nil
}

func (s *AnalyticsService) EligibleForActivity(volunteerID, activityID string) (domain.EligibilityResult, error) {
	volunteer, err := s.catalog.FindVolunteer(volunteerID)
	if err != nil {
		return domain.EligibilityResult{}, err
	}
	activity, err := s.catalog.FindActivity(activityID)
	if err != nil {
		return domain.EligibilityResult{}, err
	}
	return domain.EvaluateVolunteerEligibility(volunteer, activity, domain.DefaultEligibilityPolicy()), nil
}

func (s *AnalyticsService) RangeSummary(start, end string) (query.Result, error) {
	return s.queries.Search(query.Filter{StartDate: start, EndDate: end}, domain.ParseSortSpec("date", false))
}
