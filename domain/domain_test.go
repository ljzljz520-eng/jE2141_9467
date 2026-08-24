package domain

import "testing"

func TestVolunteerValidationAndNormalization(t *testing.T) {
	volunteer := NormalizeVolunteer(NewVolunteer(" v1 ", " Ada ", "ADA@EXAMPLE.COM ", " 1 ", "2026-01-01"))
	if err := ValidateVolunteer(volunteer); err != nil {
		t.Fatal(err)
	}
	if volunteer.Email != "ada@example.com" || volunteer.DisplayName() != "Ada" {
		t.Fatalf("unexpected volunteer: %+v", volunteer)
	}
}

func TestActivityAndRecordPolicy(t *testing.T) {
	activity := NormalizeActivity(NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 4))
	volunteer := NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")
	record := NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 3)
	if err := ValidateRecordPolicy(record, volunteer, activity, DefaultEligibilityPolicy()); err != nil {
		t.Fatal(err)
	}
	activity.Status = ActivityClosed
	if err := ValidateRecordPolicy(record, volunteer, activity, DefaultEligibilityPolicy()); err == nil {
		t.Fatal("expected closed activity error")
	}
}

func TestSortRecordsByHours(t *testing.T) {
	records := []ServiceRecord{NewServiceRecord("r2", "v1", "a1", "2026-01-02", "c", 2), NewServiceRecord("r1", "v1", "a1", "2026-01-01", "c", 5)}
	ordered := SortRecords(records, ParseSortSpec("hours", false), map[string]Volunteer{}, map[string]Activity{})
	if ordered[0].ID != "r2" {
		t.Fatalf("unexpected order: %+v", ordered)
	}
}
