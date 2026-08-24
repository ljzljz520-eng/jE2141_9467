package workflow

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
	"volunteerhours/service"
	"volunteerhours/storage"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hours.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	catalog := service.NewCatalogService(store)
	records := service.NewRecordService(store, catalog)
	if _, err := catalog.RegisterVolunteer(domain.NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")); err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.CreateActivity(domain.NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 4)); err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.OpenActivity("a1"); err != nil {
		t.Fatal(err)
	}
	if _, err := records.AppendRecord(domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 3)); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	recovery := NewRecovery(path)
	reopened, count, err := recovery.ReopenAndInspect()
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if count != 1 {
		t.Fatalf("count=%d", count)
	}
	if err := recovery.VerifyEntities(reopened, []string{"r1"}); err != nil {
		t.Fatal(err)
	}
}

func TestRecoverySnapshot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hours.db")
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SaveServiceRecord(domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "coordinator", 1)); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	recovery := NewRecovery(path)
	reopened, err := recovery.ReopenWithValidation()
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	snapshot, err := recovery.ExportSnapshot(reopened)
	if err != nil || len(snapshot) != 1 {
		t.Fatalf("snapshot=%v err=%v", snapshot, err)
	}
}
