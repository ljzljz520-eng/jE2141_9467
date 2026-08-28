package service

import (
	"volunteerhours/domain"
	"volunteerhours/query"
)

type QueryService struct {
	catalog *CatalogService
	records *RecordService
}

func NewQueryService(catalog *CatalogService, records *RecordService) *QueryService {
	return &QueryService{catalog: catalog, records: records}
}

func (s *QueryService) Search(filter query.Filter, sortSpec domain.SortSpec) (query.Result, error) {
	records, err := s.records.ListRecords()
	if err != nil {
		return query.Result{}, err
	}
	volunteers, err := s.catalog.ListVolunteers()
	if err != nil {
		return query.Result{}, err
	}
	activities, err := s.catalog.ListActivities()
	if err != nil {
		return query.Result{}, err
	}
	volunteerMap := make(map[string]domain.Volunteer, len(volunteers))
	for _, volunteer := range volunteers {
		volunteerMap[volunteer.ID] = volunteer
	}
	activityMap := make(map[string]domain.Activity, len(activities))
	for _, activity := range activities {
		activityMap[activity.ID] = activity
	}
	filtered := query.ApplyFilter(records, filter)
	sorted := domain.SortRecords(filtered, sortSpec, volunteerMap, activityMap)
	return query.Result{Records: sorted, Volunteers: volunteerMap, Activities: activityMap}, nil
}

func (s *QueryService) ByVolunteer(volunteerID string) (query.Result, error) {
	return s.Search(query.Filter{VolunteerID: volunteerID}, domain.ParseSortSpec("date", false))
}

func (s *QueryService) ByActivity(activityID string) (query.Result, error) {
	return s.Search(query.Filter{ActivityID: activityID}, domain.ParseSortSpec("date", false))
}

func (s *QueryService) Pending() (query.Result, error) {
	return s.Search(query.Filter{Status: domain.RecordPending}, domain.ParseSortSpec("date", false))
}
