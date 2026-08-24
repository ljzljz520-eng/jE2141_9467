package workflow

import (
	"fmt"

	"volunteerhours/domain"
	"volunteerhours/service"
)

type Intake struct {
	Catalog *service.CatalogService
	Records *service.RecordService
}

func NewIntake(catalog *service.CatalogService, records *service.RecordService) Intake {
	return Intake{Catalog: catalog, Records: records}
}

func (w Intake) RegisterVolunteerAndActivity(volunteer domain.Volunteer, activity domain.Activity) error {
	if _, err := w.Catalog.RegisterVolunteer(volunteer); err != nil {
		return fmt.Errorf("register volunteer: %w", err)
	}
	if _, err := w.Catalog.CreateActivity(activity); err != nil {
		return fmt.Errorf("create activity: %w", err)
	}
	if _, err := w.Catalog.OpenActivity(activity.ID); err != nil {
		return fmt.Errorf("open activity: %w", err)
	}
	return nil
}

func (w Intake) AddServiceHours(record domain.ServiceRecord) (domain.ServiceRecord, error) {
	return w.Records.AppendRecord(record)
}

func (w Intake) AddMany(records []domain.ServiceRecord) ([]domain.ServiceRecord, error) {
	added := make([]domain.ServiceRecord, 0, len(records))
	for _, record := range records {
		stored, err := w.AddServiceHours(record)
		if err != nil {
			return added, err
		}
		added = append(added, stored)
	}
	return added, nil
}

func (w Intake) AmendHours(record domain.ServiceRecord) (domain.ServiceRecord, error) {
	return w.Records.AmendRecord(record)
}
