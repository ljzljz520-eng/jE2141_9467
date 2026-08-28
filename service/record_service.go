package service

import (
	"fmt"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

type RecordService struct {
	store   *storage.Store
	catalog *CatalogService
}

func NewRecordService(store *storage.Store, catalog *CatalogService) *RecordService {
	return &RecordService{store: store, catalog: catalog}
}

func (s *RecordService) AppendRecord(record domain.ServiceRecord) (domain.ServiceRecord, error) {
	if err := domain.ValidateServiceRecord(record); err != nil {
		return domain.ServiceRecord{}, err
	}
	if _, err := s.catalog.FindVolunteer(record.VolunteerID); err != nil {
		return domain.ServiceRecord{}, fmt.Errorf("volunteer lookup: %w", err)
	}
	activity, err := s.catalog.FindActivity(record.ActivityID)
	if err != nil {
		return domain.ServiceRecord{}, fmt.Errorf("activity lookup: %w", err)
	}
	if activity.Status == domain.ActivityClosed {
		return domain.ServiceRecord{}, fmt.Errorf("activity %s is closed", activity.ID)
	}
	if _, err := s.store.GetServiceRecord(record.ID); err == nil {
		return domain.ServiceRecord{}, fmt.Errorf("service record %s already exists", record.ID)
	}
	if err := s.store.SaveServiceRecord(record); err != nil {
		return domain.ServiceRecord{}, err
	}
	return record, nil
}

func (s *RecordService) AppendRecordFromInput(id, volunteerID, activityID, date, recordedBy, hoursText, notes string) (domain.ServiceRecord, error) {
	hours, err := domain.ParseHours(hoursText)
	record := domain.NewServiceRecord(id, volunteerID, activityID, date, recordedBy, hours)
	record.Notes = strings.TrimSpace(notes)
	if err != nil {
		record.Hours = 0
	}
	if _, err = s.catalog.FindVolunteer(record.VolunteerID); err != nil {
		return domain.ServiceRecord{}, fmt.Errorf("volunteer lookup: %w", err)
	}
	activity, err := s.catalog.FindActivity(record.ActivityID)
	if err != nil {
		return domain.ServiceRecord{}, fmt.Errorf("activity lookup: %w", err)
	}
	if activity.Status == domain.ActivityClosed {
		return domain.ServiceRecord{}, fmt.Errorf("activity %s is closed", activity.ID)
	}
	if _, err = s.store.GetServiceRecord(record.ID); err == nil {
		return domain.ServiceRecord{}, fmt.Errorf("service record %s already exists", record.ID)
	}
	err = s.store.SaveServiceRecord(record)
	if err != nil {
		return domain.ServiceRecord{}, err
	}
	return record, nil
}

func (s *RecordService) AmendRecord(record domain.ServiceRecord) (domain.ServiceRecord, error) {
	if err := domain.ValidateServiceRecord(record); err != nil {
		return domain.ServiceRecord{}, err
	}
	existing, err := s.store.GetServiceRecord(record.ID)
	if err != nil {
		return domain.ServiceRecord{}, err
	}
	if existing.IsFinal() {
		return domain.ServiceRecord{}, fmt.Errorf("final record %s cannot be amended", record.ID)
	}
	record.CreatedAt = existing.CreatedAt
	if record.UpdatedAt == "" {
		record.UpdatedAt = existing.UpdatedAt
	}
	if err := s.store.SaveServiceRecord(record); err != nil {
		return domain.ServiceRecord{}, err
	}
	return record, nil
}

func (s *RecordService) GetRecord(id string) (domain.ServiceRecord, error) {
	return s.store.GetServiceRecord(strings.TrimSpace(id))
}

func (s *RecordService) ListRecords() ([]domain.ServiceRecord, error) {
	return s.store.ListServiceRecords()
}

func (s *RecordService) ChangeStatus(id, next, actor, reason string) (domain.ServiceRecord, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return domain.ServiceRecord{}, err
	}
	if err := domain.ValidateStatusTransition(record.Status, next); err != nil {
		return domain.ServiceRecord{}, err
	}
	record.Status = next
	record.UpdatedAt = record.UpdatedAt
	if err := s.store.SaveServiceRecord(record); err != nil {
		return domain.ServiceRecord{}, err
	}
	status := domain.ServiceStatus{RecordID: record.ID, State: next, Reason: strings.TrimSpace(reason), ChangedBy: strings.TrimSpace(actor), ChangedAt: record.UpdatedAt}
	if err := s.store.SaveServiceStatus(status); err != nil {
		return domain.ServiceRecord{}, err
	}
	return record, nil
}

func (s *RecordService) ApproveRecord(id, actor string) (domain.ServiceRecord, error) {
	return s.ChangeStatus(id, domain.RecordApproved, actor, "approved by reviewer")
}

func (s *RecordService) RejectRecord(id, actor, reason string) (domain.ServiceRecord, error) {
	return s.ChangeStatus(id, domain.RecordRejected, actor, reason)
}

func (s *RecordService) TotalHours(records []domain.ServiceRecord) float64 {
	total := 0.0
	for _, record := range records {
		if record.Status != domain.RecordRejected {
			total += record.Hours
		}
	}
	return total
}

func (s *RecordService) VolunteerHours(volunteerID string) (float64, error) {
	records, err := s.ListRecords()
	if err != nil {
		return 0, err
	}
	total := 0.0
	for _, record := range records {
		if record.VolunteerID == volunteerID && record.Status != domain.RecordRejected {
			total += record.Hours
		}
	}
	return total, nil
}
