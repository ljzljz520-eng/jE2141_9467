package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type ValidationIssue struct {
	Field   string
	Message string
}

func (i ValidationIssue) Error() string {
	return fmt.Sprintf("%s: %s", i.Field, i.Message)
}

type ValidationErrors []ValidationIssue

func (e ValidationErrors) Error() string {
	parts := make([]string, 0, len(e))
	for _, issue := range e {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}

func (e ValidationErrors) Has(field string) bool {
	for _, issue := range e {
		if issue.Field == field {
			return true
		}
	}
	return false
}

func ValidateVolunteer(v Volunteer) error {
	issues := make(ValidationErrors, 0, 4)
	if strings.TrimSpace(v.ID) == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "is required"})
	}
	if strings.TrimSpace(v.Name) == "" {
		issues = append(issues, ValidationIssue{Field: "name", Message: "is required"})
	}
	if strings.TrimSpace(v.JoinedDate) == "" {
		issues = append(issues, ValidationIssue{Field: "joined_date", Message: "is required"})
	} else if _, err := NormalizeDate(v.JoinedDate); err != nil {
		issues = append(issues, ValidationIssue{Field: "joined_date", Message: err.Error()})
	}
	if v.Email != "" && !emailPattern.MatchString(v.Email) {
		issues = append(issues, ValidationIssue{Field: "email", Message: "is invalid"})
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func ValidateActivity(a Activity) error {
	issues := make(ValidationErrors, 0, 6)
	if strings.TrimSpace(a.ID) == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "is required"})
	}
	if strings.TrimSpace(a.Title) == "" {
		issues = append(issues, ValidationIssue{Field: "title", Message: "is required"})
	}
	if a.Capacity < 0 {
		issues = append(issues, ValidationIssue{Field: "capacity", Message: "cannot be negative"})
	}
	if a.StartDate == "" {
		issues = append(issues, ValidationIssue{Field: "start_date", Message: "is required"})
	} else if _, err := NormalizeDate(a.StartDate); err != nil {
		issues = append(issues, ValidationIssue{Field: "start_date", Message: err.Error()})
	}
	if a.EndDate == "" {
		issues = append(issues, ValidationIssue{Field: "end_date", Message: "is required"})
	} else if _, err := NormalizeDate(a.EndDate); err != nil {
		issues = append(issues, ValidationIssue{Field: "end_date", Message: err.Error()})
	}
	if a.StartDate != "" && a.EndDate != "" && a.StartDate > a.EndDate {
		issues = append(issues, ValidationIssue{Field: "end_date", Message: "must not precede start_date"})
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func ValidateServiceRecord(r ServiceRecord) error {
	issues := make(ValidationErrors, 0, 8)
	if strings.TrimSpace(r.ID) == "" {
		issues = append(issues, ValidationIssue{Field: "id", Message: "is required"})
	}
	if strings.TrimSpace(r.VolunteerID) == "" {
		issues = append(issues, ValidationIssue{Field: "volunteer_id", Message: "is required"})
	}
	if strings.TrimSpace(r.ActivityID) == "" {
		issues = append(issues, ValidationIssue{Field: "activity_id", Message: "is required"})
	}
	if _, err := NormalizeDate(r.ServiceDate); err != nil {
		issues = append(issues, ValidationIssue{Field: "service_date", Message: err.Error()})
	}
	if r.Hours <= 0 {
		issues = append(issues, ValidationIssue{Field: "hours", Message: "must be greater than zero"})
	} else if r.Hours > 24 {
		issues = append(issues, ValidationIssue{Field: "hours", Message: "must not exceed 24"})
	}
	if strings.TrimSpace(r.RecordedBy) == "" {
		issues = append(issues, ValidationIssue{Field: "recorded_by", Message: "is required"})
	}
	if len(issues) > 0 {
		return issues
	}
	return nil
}

func ValidateStatusTransition(current, next string) error {
	if current == next {
		return nil
	}
	if current == RecordPending && (next == RecordApproved || next == RecordRejected) {
		return nil
	}
	return fmt.Errorf("cannot change status from %s to %s", current, next)
}

func ValidateDateRange(start, end string) error {
	if strings.TrimSpace(start) == "" || strings.TrimSpace(end) == "" {
		return fmt.Errorf("both dates are required")
	}
	if start > end {
		return fmt.Errorf("start date must not follow end date")
	}
	return nil
}
