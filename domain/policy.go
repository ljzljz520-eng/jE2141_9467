package domain

import (
	"fmt"
	"strings"
)

type EligibilityPolicy struct {
	MinimumAgeLabel string
	MaximumHours    float64
	RequiresActive  bool
	AllowedStatuses []string
}

type EligibilityResult struct {
	Eligible bool
	Reasons  []string
}

func DefaultEligibilityPolicy() EligibilityPolicy {
	return EligibilityPolicy{MinimumAgeLabel: "registered", MaximumHours: 24, RequiresActive: true, AllowedStatuses: []string{RecordPending, RecordApproved}}
}

func EvaluateVolunteerEligibility(volunteer Volunteer, activity Activity, policy EligibilityPolicy) EligibilityResult {
	result := EligibilityResult{Eligible: true, Reasons: []string{}}
	if !volunteer.IsComplete() {
		result.Eligible = false
		result.Reasons = append(result.Reasons, "volunteer profile is incomplete")
	}
	if policy.RequiresActive && !volunteer.Active {
		result.Eligible = false
		result.Reasons = append(result.Reasons, "volunteer is inactive")
	}
	if activity.Status == ActivityClosed {
		result.Eligible = false
		result.Reasons = append(result.Reasons, "activity is closed")
	}
	if activity.Capacity == 0 {
		result.Eligible = false
		result.Reasons = append(result.Reasons, "activity has no capacity")
	}
	return result
}

func ValidateRecordPolicy(record ServiceRecord, volunteer Volunteer, activity Activity, policy EligibilityPolicy) error {
	if err := ValidateServiceRecord(record); err != nil {
		return err
	}
	eligibility := EvaluateVolunteerEligibility(volunteer, activity, policy)
	if !eligibility.Eligible {
		return fmt.Errorf("record policy rejected: %s", strings.Join(eligibility.Reasons, ", "))
	}
	if policy.MaximumHours > 0 && record.Hours > policy.MaximumHours {
		return fmt.Errorf("record policy limit is %.2f hours", policy.MaximumHours)
	}
	return nil
}

func IsAllowedStatus(status string, allowed []string) bool {
	for _, item := range allowed {
		if status == item {
			return true
		}
	}
	return false
}

func StatusDescription(status string) string {
	switch status {
	case RecordPending:
		return "awaiting coordinator review"
	case RecordApproved:
		return "included in official totals"
	case RecordRejected:
		return "excluded from official totals"
	default:
		return "unknown status"
	}
}

func NormalizeVolunteer(v Volunteer) Volunteer {
	v.ID = strings.TrimSpace(v.ID)
	v.Name = strings.TrimSpace(v.Name)
	v.Email = strings.ToLower(strings.TrimSpace(v.Email))
	v.Phone = strings.TrimSpace(v.Phone)
	v.JoinedDate = strings.TrimSpace(v.JoinedDate)
	v.Notes = strings.TrimSpace(v.Notes)
	return v
}

func NormalizeActivity(a Activity) Activity {
	a.ID = strings.TrimSpace(a.ID)
	a.Title = strings.TrimSpace(a.Title)
	a.Location = strings.TrimSpace(a.Location)
	a.Coordinator = strings.TrimSpace(a.Coordinator)
	a.StartDate = strings.TrimSpace(a.StartDate)
	a.EndDate = strings.TrimSpace(a.EndDate)
	a.Status = strings.ToLower(strings.TrimSpace(a.Status))
	a.Description = strings.TrimSpace(a.Description)
	if a.Status == "" {
		a.Status = ActivityDraft
	}
	return a
}

func NormalizeServiceRecord(r ServiceRecord) ServiceRecord {
	r.ID = strings.TrimSpace(r.ID)
	r.VolunteerID = strings.TrimSpace(r.VolunteerID)
	r.ActivityID = strings.TrimSpace(r.ActivityID)
	r.ServiceDate = strings.TrimSpace(r.ServiceDate)
	r.RecordedBy = strings.TrimSpace(r.RecordedBy)
	r.Status = strings.ToLower(strings.TrimSpace(r.Status))
	r.Notes = strings.TrimSpace(r.Notes)
	if r.Status == "" {
		r.Status = RecordPending
	}
	return r
}

func MergeVolunteer(existing, update Volunteer) Volunteer {
	merged := existing
	if strings.TrimSpace(update.Name) != "" {
		merged.Name = strings.TrimSpace(update.Name)
	}
	if strings.TrimSpace(update.Email) != "" {
		merged.Email = strings.TrimSpace(update.Email)
	}
	if strings.TrimSpace(update.Phone) != "" {
		merged.Phone = strings.TrimSpace(update.Phone)
	}
	if strings.TrimSpace(update.Notes) != "" {
		merged.Notes = strings.TrimSpace(update.Notes)
	}
	merged.Active = update.Active
	return merged
}

func MergeActivity(existing, update Activity) Activity {
	merged := existing
	if strings.TrimSpace(update.Title) != "" {
		merged.Title = strings.TrimSpace(update.Title)
	}
	if strings.TrimSpace(update.Location) != "" {
		merged.Location = strings.TrimSpace(update.Location)
	}
	if strings.TrimSpace(update.Coordinator) != "" {
		merged.Coordinator = strings.TrimSpace(update.Coordinator)
	}
	if strings.TrimSpace(update.StartDate) != "" {
		merged.StartDate = strings.TrimSpace(update.StartDate)
	}
	if strings.TrimSpace(update.EndDate) != "" {
		merged.EndDate = strings.TrimSpace(update.EndDate)
	}
	if update.Capacity > 0 {
		merged.Capacity = update.Capacity
	}
	if strings.TrimSpace(update.Description) != "" {
		merged.Description = strings.TrimSpace(update.Description)
	}
	return merged
}

func RecordKey(record ServiceRecord) string {
	return strings.Join([]string{record.VolunteerID, record.ActivityID, record.ServiceDate}, ":")
}

func CompareRecordIdentity(left, right ServiceRecord) bool {
	return left.ID == right.ID && left.VolunteerID == right.VolunteerID && left.ActivityID == right.ActivityID && left.ServiceDate == right.ServiceDate
}

func IsDateWithin(value, start, end string) bool {
	if start != "" && value < start {
		return false
	}
	if end != "" && value > end {
		return false
	}
	return true
}
