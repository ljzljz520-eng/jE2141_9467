package storage

import (
	"fmt"
	"sort"

	"volunteerhours/domain"
)

func (s *Store) SaveServiceRecord(record domain.ServiceRecord) error {
	data, err := domain.Encode(record)
	if err != nil {
		return fmt.Errorf("encode ServiceRecord: %w", err)
	}
	return s.write(recordBucket, record.ID, data)
}

func (s *Store) GetServiceRecord(id string) (domain.ServiceRecord, error) {
	data, err := s.read(recordBucket, id)
	if err != nil {
		return domain.ServiceRecord{}, err
	}
	var record domain.ServiceRecord
	if err := domain.Decode(data, &record); err != nil {
		return domain.ServiceRecord{}, fmt.Errorf("decode ServiceRecord: %w", err)
	}
	return record, nil
}

func (s *Store) ListServiceRecords() ([]domain.ServiceRecord, error) {
	items, err := s.list(recordBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ServiceRecord, 0, len(items))
	for _, data := range items {
		var record domain.ServiceRecord
		if err := domain.Decode(data, &record); err != nil {
			return nil, fmt.Errorf("decode ServiceRecord list: %w", err)
		}
		result = append(result, record)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *Store) DeleteServiceRecord(id string) error {
	return s.remove(recordBucket, id)
}

func (s *Store) SaveServiceStatus(status domain.ServiceStatus) error {
	data, err := domain.Encode(status)
	if err != nil {
		return fmt.Errorf("encode ServiceStatus: %w", err)
	}
	return s.write(statusBucket, status.RecordID, data)
}

func (s *Store) GetServiceStatus(id string) (domain.ServiceStatus, error) {
	data, err := s.read(statusBucket, id)
	if err != nil {
		return domain.ServiceStatus{}, err
	}
	var status domain.ServiceStatus
	if err := domain.Decode(data, &status); err != nil {
		return domain.ServiceStatus{}, fmt.Errorf("decode ServiceStatus: %w", err)
	}
	return status, nil
}

func (s *Store) SaveAuditEntry(entry domain.AuditEntry) error {
	data, err := domain.Encode(entry)
	if err != nil {
		return fmt.Errorf("encode AuditEntry: %w", err)
	}
	return s.write(auditBucket, entry.ID, data)
}

func (s *Store) ListAuditEntries() ([]domain.AuditEntry, error) {
	items, err := s.list(auditBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEntry, 0, len(items))
	for _, data := range items {
		var entry domain.AuditEntry
		if err := domain.Decode(data, &entry); err != nil {
			return nil, fmt.Errorf("decode AuditEntry list: %w", err)
		}
		result = append(result, entry)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *Store) RecordCount() (int, error) {
	return s.Count(recordBucket)
}
