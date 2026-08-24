package service

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

func recordFixture(t *testing.T) (*storage.Store, *CatalogService, *RecordService) {
	t.Helper()
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	catalog := NewCatalogService(store)
	records := NewRecordService(store, catalog)
	if _, err := catalog.RegisterVolunteer(domain.NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")); err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.CreateActivity(domain.NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 2)); err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.OpenActivity("a1"); err != nil {
		t.Fatal(err)
	}
	return store, catalog, records
}

func TestRecordApprovalLifecycle(t *testing.T) {
	store, _, records := recordFixture(t)
	defer store.Close()
	record, err := records.AppendRecord(domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 3))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := records.ApproveRecord(record.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	got, err := records.GetRecord("r1")
	if err != nil || got.Status != domain.RecordApproved {
		t.Fatalf("record=%+v err=%v", got, err)
	}
	if _, err := records.AmendRecord(got); err == nil {
		t.Fatal("expected final record error")
	}
}

func TestRecordTotalsExcludeRejected(t *testing.T) {
	store, _, records := recordFixture(t)
	defer store.Close()
	first := domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 3)
	second := domain.NewServiceRecord("r2", "v1", "a1", "2026-02-02", "coordinator", 2)
	if _, err := records.AppendRecord(first); err != nil {
		t.Fatal(err)
	}
	if _, err := records.AppendRecord(second); err != nil {
		t.Fatal(err)
	}
	if _, err := records.RejectRecord("r2", "reviewer", "duplicate"); err != nil {
		t.Fatal(err)
	}
	hours, err := records.VolunteerHours("v1")
	if err != nil || hours != 3 {
		t.Fatalf("hours=%v err=%v", hours, err)
	}
}
