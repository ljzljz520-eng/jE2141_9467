package workflow

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
	"volunteerhours/service"
	"volunteerhours/storage"
)

func TestWorkflowOne(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	catalog := service.NewCatalogService(store)
	records := service.NewRecordService(store, catalog)
	intake := NewIntake(catalog, records)
	volunteer := domain.NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")
	activity := domain.NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 4)
	if err := intake.RegisterVolunteerAndActivity(volunteer, activity); err != nil {
		t.Fatal(err)
	}
	if _, err := intake.AddServiceHours(domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 3)); err != nil {
		t.Fatal(err)
	}
	if got, err := records.ListRecords(); err != nil || len(got) != 1 {
		t.Fatalf("records=%v err=%v", got, err)
	}
}
