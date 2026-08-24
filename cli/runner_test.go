package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"volunteerhours/storage"
)

func TestWorkflowThree(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	var output bytes.Buffer
	runner := NewRunner(strings.NewReader(""), &output, store)
	if err := runner.Run([]string{"volunteer", "add", "-id", "v1", "-name", "Ada", "-email", "ada@example.com", "-joined", "2026-01-01"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"activity", "add", "-id", "a1", "-title", "Cleanup", "-location", "Park", "-coordinator", "Ada", "-start", "2026-02-01", "-end", "2026-02-01", "-capacity", "4"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"activity", "open", "a1"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"record", "add", "-id", "r1", "-volunteer", "v1", "-activity", "a1", "-date", "2026-02-01", "-hours", "3", "-by", "coordinator"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"review", "approve", "r1"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"summary"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Total hours: 3.00") {
		t.Fatalf("output=%s", output.String())
	}
}

func TestInvalidVolunteerHoursAreRejected(t *testing.T) {
	store, err := storage.Open(filepath.Join(t.TempDir(), "hours.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	runner := NewRunner(strings.NewReader(""), &bytes.Buffer{}, store)
	if err := runner.Run([]string{"volunteer", "add", "-id", "v1", "-name", "Ada", "-joined", "2026-01-01"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"activity", "add", "-id", "a1", "-title", "Cleanup", "-start", "2026-02-01", "-end", "2026-02-01", "-capacity", "4"}); err != nil {
		t.Fatal(err)
	}
	if err := runner.Run([]string{"activity", "open", "a1"}); err != nil {
		t.Fatal(err)
	}
	err = runner.Run([]string{"record", "add", "-id", "bad-hours", "-volunteer", "v1", "-activity", "a1", "-date", "2026-02-01", "-hours", "three", "-by", "coordinator"})
	if err == nil {
		t.Fatal("expected invalid hours error")
	}
	count, countErr := store.RecordCount()
	if countErr != nil {
		t.Fatal(countErr)
	}
	if count != 0 {
		t.Fatalf("invalid record count=%d", count)
	}
}
