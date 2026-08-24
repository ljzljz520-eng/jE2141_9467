package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Volunteer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	JoinedDate string `json:"joined_date"`
	Active     bool   `json:"active"`
	Notes      string `json:"notes"`
}

type Activity struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Location    string `json:"location"`
	Coordinator string `json:"coordinator"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Capacity    int    `json:"capacity"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ServiceRecord struct {
	ID          string  `json:"id"`
	VolunteerID string  `json:"volunteer_id"`
	ActivityID  string  `json:"activity_id"`
	ServiceDate string  `json:"service_date"`
	Hours       float64 `json:"hours"`
	RecordedBy  string  `json:"recorded_by"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ServiceStatus struct {
	RecordID  string `json:"record_id"`
	State     string `json:"state"`
	Reason    string `json:"reason"`
	ChangedBy string `json:"changed_by"`
	ChangedAt string `json:"changed_at"`
}

type AuditEntry struct {
	ID         string `json:"id"`
	Entity     string `json:"entity"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	OccurredAt string `json:"occurred_at"`
	Detail     string `json:"detail"`
}

const (
	RecordPending  = "pending"
	RecordApproved = "approved"
	RecordRejected = "rejected"
	ActivityDraft  = "draft"
	ActivityOpen   = "open"
	ActivityClosed = "closed"
)

func NewVolunteer(id, name, email, phone, joined string) Volunteer {
	return Volunteer{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name), Email: strings.TrimSpace(email), Phone: strings.TrimSpace(phone), JoinedDate: strings.TrimSpace(joined), Active: true}
}

func NewActivity(id, title, location, coordinator, start, end string, capacity int) Activity {
	return Activity{ID: strings.TrimSpace(id), Title: strings.TrimSpace(title), Location: strings.TrimSpace(location), Coordinator: strings.TrimSpace(coordinator), StartDate: strings.TrimSpace(start), EndDate: strings.TrimSpace(end), Capacity: capacity, Status: ActivityDraft}
}

func NewServiceRecord(id, volunteerID, activityID, date, recordedBy string, hours float64) ServiceRecord {
	return ServiceRecord{ID: strings.TrimSpace(id), VolunteerID: strings.TrimSpace(volunteerID), ActivityID: strings.TrimSpace(activityID), ServiceDate: strings.TrimSpace(date), Hours: hours, RecordedBy: strings.TrimSpace(recordedBy), Status: RecordPending, CreatedAt: "fixed", UpdatedAt: "fixed"}
}

func (v Volunteer) IsComplete() bool {
	return strings.TrimSpace(v.ID) != "" && strings.TrimSpace(v.Name) != "" && strings.TrimSpace(v.JoinedDate) != ""
}

func (v Volunteer) DisplayName() string {
	if strings.TrimSpace(v.Name) == "" {
		return v.ID
	}
	return v.Name
}

func (a Activity) IsScheduled() bool {
	return a.StartDate != "" && a.EndDate != "" && a.Status != ActivityClosed
}

func (a Activity) Label() string {
	if a.Location == "" {
		return a.Title
	}
	return fmt.Sprintf("%s at %s", a.Title, a.Location)
}

func (r ServiceRecord) IsValid() bool {
	return strings.TrimSpace(r.ID) != "" && strings.TrimSpace(r.VolunteerID) != "" && strings.TrimSpace(r.ActivityID) != "" && strings.TrimSpace(r.ServiceDate) != "" && r.Hours > 0 && r.Hours <= 24 && strings.TrimSpace(r.RecordedBy) != ""
}

func (r ServiceRecord) IsFinal() bool {
	return r.Status == RecordApproved || r.Status == RecordRejected
}

func (r ServiceRecord) HoursLabel() string {
	return fmt.Sprintf("%.2f hours", r.Hours)
}

func (s ServiceStatus) IsTerminal() bool {
	return s.State == RecordApproved || s.State == RecordRejected
}

func Encode[T any](value T) ([]byte, error) {
	return json.Marshal(value)
}

func Decode[T any](data []byte, target *T) error {
	return json.Unmarshal(data, target)
}

func NormalizeDate(value string) (string, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return "", fmt.Errorf("date must use YYYY-MM-DD: %w", err)
	}
	return parsed.Format("2006-01-02"), nil
}

func ParseHours(value string) (float64, error) {
	var hours float64
	if _, err := fmt.Sscanf(strings.TrimSpace(value), "%f", &hours); err != nil {
		return 0, fmt.Errorf("hours must be numeric: %w", err)
	}
	return hours, nil
}
