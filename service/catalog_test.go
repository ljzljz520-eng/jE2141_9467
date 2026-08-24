package service

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

func TestCatalogLifecycle(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	catalog := NewCatalogService(store)
	volunteer := domain.NewVolunteer("v1", "Ada", "ada@example.com", "", "2026-01-01")
	if _, err := catalog.RegisterVolunteer(volunteer); err != nil {
		t.Fatal(err)
	}
	activity := domain.NewActivity("a1", "Cleanup", "Park", "Ada", "2026-02-01", "2026-02-01", 2)
	if _, err := catalog.CreateActivity(activity); err != nil {
		t.Fatal(err)
	}
	if got, err := catalog.OpenActivity("a1"); err != nil || got.Status != domain.ActivityOpen {
		t.Fatalf("activity=%+v err=%v", got, err)
	}
	if _, err := catalog.CloseActivity("a1"); err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.RegisterVolunteer(volunteer); err == nil {
		t.Fatal("expected duplicate volunteer error")
	}
}
