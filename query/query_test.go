package query

import (
	"testing"

	"volunteerhours/domain"
)

func TestFilterGroupAndAggregate(t *testing.T) {
	records := []domain.ServiceRecord{
		domain.NewServiceRecord("r1", "v1", "a1", "2026-02-01", "Ada", 3),
		domain.NewServiceRecord("r2", "v1", "a2", "2026-02-02", "Ada", 2),
		domain.NewServiceRecord("r3", "v2", "a1", "2026-02-02", "Ben", 4),
	}
	filtered := ApplyFilter(records, Filter{VolunteerID: "v1", MinHours: 2})
	if len(filtered) != 2 {
		t.Fatalf("filtered=%d", len(filtered))
	}
	result := Result{Records: filtered, Volunteers: map[string]domain.Volunteer{"v1": {ID: "v1", Name: "Ada"}}, Activities: map[string]domain.Activity{"a1": {ID: "a1", Title: "Cleanup"}, "a2": {ID: "a2", Title: "Food"}}}
	if groups := GroupByActivity(result); len(groups) != 2 {
		t.Fatalf("groups=%v", groups)
	}
	if got := DailyTotals(records); len(got) != 2 || got[0].Hours != 3 {
		t.Fatalf("daily=%v", got)
	}
}
