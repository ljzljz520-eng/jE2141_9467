package storage

import (
	"fmt"
	"sort"
	"strings"

	"volunteerhours/domain"
)

func (s *Store) FindRecordsByVolunteer(volunteerID string) ([]domain.ServiceRecord, error) {
	records, err := s.ListServiceRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.ServiceRecord, 0)
	for _, record := range records {
		if record.VolunteerID == volunteerID {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Store) FindRecordsByActivity(activityID string) ([]domain.ServiceRecord, error) {
	records, err := s.ListServiceRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.ServiceRecord, 0)
	for _, record := range records {
		if record.ActivityID == activityID {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Store) FindRecordsInRange(start, end string) ([]domain.ServiceRecord, error) {
	if err := domain.ValidateDateRange(start, end); err != nil {
		return nil, err
	}
	records, err := s.ListServiceRecords()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.ServiceRecord, 0)
	for _, record := range records {
		if domain.IsDateWithin(record.ServiceDate, start, end) {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Store) FindDuplicateServiceRecord(candidate domain.ServiceRecord) (domain.ServiceRecord, bool, error) {
	records, err := s.ListServiceRecords()
	if err != nil {
		return domain.ServiceRecord{}, false, err
	}
	for _, record := range records {
		if domain.RecordKey(record) == domain.RecordKey(candidate) {
			return record, true, nil
		}
	}
	return domain.ServiceRecord{}, false, nil
}

func (s *Store) SaveAudit(actor, entity, entityID, action, detail string) (domain.AuditEntry, error) {
	if strings.TrimSpace(actor) == "" || strings.TrimSpace(entityID) == "" {
		return domain.AuditEntry{}, fmt.Errorf("audit actor and entity id are required")
	}
	entry := domain.AuditEntry{ID: fmt.Sprintf("%s-%s-%s", entity, entityID, action), Entity: entity, EntityID: entityID, Action: action, Actor: actor, OccurredAt: "fixed", Detail: detail}
	if err := s.SaveAuditEntry(entry); err != nil {
		return domain.AuditEntry{}, err
	}
	return entry, nil
}

func (s *Store) AuditForEntity(entity, entityID string) ([]domain.AuditEntry, error) {
	entries, err := s.ListAuditEntries()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.AuditEntry, 0)
	for _, entry := range entries {
		if entry.Entity == entity && entry.EntityID == entityID {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func (s *Store) ListStatuses() ([]domain.ServiceStatus, error) {
	items, err := s.list(statusBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ServiceStatus, 0, len(items))
	for _, data := range items {
		var status domain.ServiceStatus
		if err := domain.Decode(data, &status); err != nil {
			return nil, fmt.Errorf("decode ServiceStatus list: %w", err)
		}
		result = append(result, status)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].RecordID < result[j].RecordID })
	return result, nil
}

func (s *Store) EntityCounts() (map[string]int, error) {
	volunteers, err := s.VolunteerCount()
	if err != nil {
		return nil, err
	}
	activities, err := s.ActivityCount()
	if err != nil {
		return nil, err
	}
	records, err := s.RecordCount()
	if err != nil {
		return nil, err
	}
	audits, err := s.Count(auditBucket)
	if err != nil {
		return nil, err
	}
	return map[string]int{"Volunteer": volunteers, "Activity": activities, "ServiceRecord": records, "AuditEntry": audits}, nil
}

func (s *Store) SaveVolunteerWithAudit(volunteer domain.Volunteer, actor string) error {
	if err := s.SaveVolunteer(volunteer); err != nil {
		return err
	}
	_, err := s.SaveAudit(actor, "Volunteer", volunteer.ID, "save", volunteer.DisplayName())
	return err
}

func (s *Store) SaveActivityWithAudit(activity domain.Activity, actor string) error {
	if err := s.SaveActivity(activity); err != nil {
		return err
	}
	_, err := s.SaveAudit(actor, "Activity", activity.ID, "save", activity.Label())
	return err
}

func (s *Store) SaveRecordWithAudit(record domain.ServiceRecord, actor string) error {
	if err := s.SaveServiceRecord(record); err != nil {
		return err
	}
	_, err := s.SaveAudit(actor, "ServiceRecord", record.ID, "save", record.HoursLabel())
	return err
}
