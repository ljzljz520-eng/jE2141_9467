package storage

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
)

func TestStorePersistsEntities(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	volunteer := domain.NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")
	activity := domain.NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 10)
	record := domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 2)
	if err := store.SaveVolunteer(volunteer); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveActivity(activity); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveServiceRecord(record); err != nil {
		t.Fatal(err)
	}
	if got, err := store.RecordCount(); err != nil || got != 1 {
		t.Fatalf("record count=%d err=%v", got, err)
	}
	if got, err := store.GetVolunteer("v1"); err != nil || got.Name != "Ada" {
		t.Fatalf("volunteer=%+v err=%v", got, err)
	}
}

func TestStoreQueriesAndAudit(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	record := domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 2)
	if err := store.SaveServiceRecord(record); err != nil {
		t.Fatal(err)
	}
	if _, err := store.SaveAudit("coordinator", "ServiceRecord", "r1", "save", "2 hours"); err != nil {
		t.Fatal(err)
	}
	found, ok, err := store.FindDuplicateServiceRecord(record)
	if err != nil || !ok || found.ID != "r1" {
		t.Fatalf("duplicate=%+v ok=%v err=%v", found, ok, err)
	}
	entries, err := store.AuditForEntity("ServiceRecord", "r1")
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
}
