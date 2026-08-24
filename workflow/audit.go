package workflow

import (
	"fmt"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

type AuditWorkflow struct {
	Store *storage.Store
}

func NewAuditWorkflow(store *storage.Store) AuditWorkflow {
	return AuditWorkflow{Store: store}
}

func (w AuditWorkflow) RecordVolunteerChange(volunteer domain.Volunteer, actor, action string) (domain.AuditEntry, error) {
	if strings.TrimSpace(action) == "" {
		return domain.AuditEntry{}, fmt.Errorf("audit action is required")
	}
	return w.Store.SaveAudit(actor, "Volunteer", volunteer.ID, action, volunteer.DisplayName())
}

func (w AuditWorkflow) RecordActivityChange(activity domain.Activity, actor, action string) (domain.AuditEntry, error) {
	if activity.ID == "" {
		return domain.AuditEntry{}, fmt.Errorf("activity id is required")
	}
	return w.Store.SaveAudit(actor, "Activity", activity.ID, action, activity.Status)
}

func (w AuditWorkflow) RecordServiceChange(record domain.ServiceRecord, actor, action string) (domain.AuditEntry, error) {
	if record.ID == "" {
		return domain.AuditEntry{}, fmt.Errorf("record id is required")
	}
	return w.Store.SaveAudit(actor, "ServiceRecord", record.ID, action, record.Status)
}

func (w AuditWorkflow) History(entity, id string) ([]domain.AuditEntry, error) {
	return w.Store.AuditForEntity(entity, id)
}

func (w AuditWorkflow) Describe(entries []domain.AuditEntry) string {
	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		parts = append(parts, fmt.Sprintf("%s:%s", entry.Action, entry.Actor))
	}
	return strings.Join(parts, ", ")
}
