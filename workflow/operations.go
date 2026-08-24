package workflow

import (
	"fmt"

	"volunteerhours/domain"
	"volunteerhours/query"
	"volunteerhours/service"
)

type Operations struct {
	Catalog *service.CatalogService
	Records *service.RecordService
	Query   *service.QueryService
}

func NewOperations(catalog *service.CatalogService, records *service.RecordService, queryService *service.QueryService) Operations {
	return Operations{Catalog: catalog, Records: records, Query: queryService}
}

func (w Operations) ReviewPending(id, reviewer string, approve bool, reason string) (domain.ServiceRecord, error) {
	if approve {
		return w.Records.ApproveRecord(id, reviewer)
	}
	if reason == "" {
		return domain.ServiceRecord{}, fmt.Errorf("rejection reason is required")
	}
	return w.Records.RejectRecord(id, reviewer, reason)
}

func (w Operations) Search(filter query.Filter, sortSpec domain.SortSpec, offset, limit int) (query.Result, error) {
	result, err := w.Query.Search(filter, sortSpec)
	if err != nil {
		return query.Result{}, err
	}
	return query.Limit(result, offset, limit), nil
}

func (w Operations) CloseActivity(id string) (domain.Activity, error) {
	return w.Catalog.CloseActivity(id)
}

func (w Operations) UpdateVolunteer(volunteer domain.Volunteer) (domain.Volunteer, error) {
	return w.Catalog.UpdateVolunteer(volunteer)
}

func (w Operations) UpdateActivity(activity domain.Activity) (domain.Activity, error) {
	return w.Catalog.UpdateActivity(activity)
}

func (w Operations) Summary(filter query.Filter) (query.Result, float64, error) {
	result, err := w.Query.Search(filter, domain.ParseSortSpec("date", false))
	if err != nil {
		return query.Result{}, 0, err
	}
	return result, result.Hours(), nil
}
