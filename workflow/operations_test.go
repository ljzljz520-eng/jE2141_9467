package workflow

import (
	"path/filepath"
	"testing"

	"volunteerhours/domain"
	"volunteerhours/query"
	"volunteerhours/service"
	"volunteerhours/storage"
)

func TestWorkflowTwo(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	catalog := service.NewCatalogService(store)
	records := service.NewRecordService(store, catalog)
	queries := service.NewQueryService(catalog, records)
	ops := NewOperations(catalog, records, queries)
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
	if _, err := ops.ReviewPending("r1", "reviewer", true, ""); err != nil {
		t.Fatal(err)
	}
	result, total, err := ops.Summary(query.Filter{Status: domain.RecordApproved})
	if err != nil || result.Count() != 1 || total != 3 {
		t.Fatalf("result=%v total=%v err=%v", result, total, err)
	}
}
