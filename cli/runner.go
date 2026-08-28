package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/query"
	"volunteerhours/report"
	"volunteerhours/service"
	"volunteerhours/storage"
	"volunteerhours/workflow"
)

type Runner struct {
	Input  io.Reader
	Output io.Writer
	Store  *storage.Store
	Intake workflow.Intake
	Ops    workflow.Operations
}

func NewRunner(input io.Reader, output io.Writer, store *storage.Store) *Runner {
	catalog := service.NewCatalogService(store)
	records := service.NewRecordService(store, catalog)
	queries := service.NewQueryService(catalog, records)
	return &Runner{Input: input, Output: output, Store: store, Intake: workflow.NewIntake(catalog, records), Ops: workflow.NewOperations(catalog, records, queries)}
}

func (r *Runner) Run(args []string) error {
	if len(args) == 0 {
		return r.printUsage()
	}
	switch args[0] {
	case "volunteer":
		return r.runVolunteer(args[1:])
	case "activity":
		return r.runActivity(args[1:])
	case "record":
		return r.runRecord(args[1:])
	case "review":
		return r.runReview(args[1:])
	case "list":
		return r.runList(args[1:])
	case "summary":
		return r.runSummary(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (r *Runner) printUsage() error {
	_, err := fmt.Fprintln(r.Output, "commands: volunteer add|list, activity add|open|close|list, record add|list, review approve|reject, list records, summary")
	return err
}

func (r *Runner) runVolunteer(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("volunteer action is required")
	}
	switch args[0] {
	case "add":
		flags := flag.NewFlagSet("volunteer add", flag.ContinueOnError)
		flags.SetOutput(io.Discard)
		id := flags.String("id", "", "id")
		name := flags.String("name", "", "name")
		email := flags.String("email", "", "email")
		phone := flags.String("phone", "", "phone")
		joined := flags.String("joined", "", "joined date")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		volunteer := domain.NewVolunteer(*id, *name, *email, *phone, *joined)
		stored, err := r.Intake.Catalog.RegisterVolunteer(volunteer)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(r.Output, "saved volunteer %s\n", stored.ID)
		return err
	case "list":
		volunteers, err := r.Intake.Catalog.ListVolunteers()
		if err != nil {
			return err
		}
		return report.WriteVolunteers(r.Output, volunteers)
	default:
		return fmt.Errorf("unknown volunteer action %q", args[0])
	}
}

func (r *Runner) runActivity(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("activity action is required")
	}
	switch args[0] {
	case "add":
		flags := flag.NewFlagSet("activity add", flag.ContinueOnError)
		flags.SetOutput(io.Discard)
		id := flags.String("id", "", "id")
		title := flags.String("title", "", "title")
		location := flags.String("location", "", "location")
		coordinator := flags.String("coordinator", "", "coordinator")
		start := flags.String("start", "", "start")
		end := flags.String("end", "", "end")
		capacity := flags.Int("capacity", 0, "capacity")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		activity := domain.NewActivity(*id, *title, *location, *coordinator, *start, *end, *capacity)
		stored, err := r.Intake.Catalog.CreateActivity(activity)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(r.Output, "saved activity %s\n", stored.ID)
		return err
	case "open", "close":
		if len(args) < 2 {
			return fmt.Errorf("activity id is required")
		}
		var activity domain.Activity
		var err error
		if args[0] == "open" {
			activity, err = r.Intake.Catalog.OpenActivity(args[1])
		} else {
			activity, err = r.Intake.Catalog.CloseActivity(args[1])
		}
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(r.Output, "%s activity %s\n", activity.Status, activity.ID)
		return err
	case "list":
		activities, err := r.Intake.Catalog.ListActivities()
		if err != nil {
			return err
		}
		return report.WriteActivities(r.Output, activities)
	default:
		return fmt.Errorf("unknown activity action %q", args[0])
	}
}

func (r *Runner) runRecord(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("record action is required")
	}
	switch args[0] {
	case "add":
		flags := flag.NewFlagSet("record add", flag.ContinueOnError)
		flags.SetOutput(io.Discard)
		id := flags.String("id", "", "id")
		volunteerID := flags.String("volunteer", "", "volunteer")
		activityID := flags.String("activity", "", "activity")
		date := flags.String("date", "", "date")
		hoursText := flags.String("hours", "", "hours")
		recordedBy := flags.String("by", "", "recorded by")
		notes := flags.String("notes", "", "notes")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		stored, err := r.Intake.Records.AppendRecordFromInput(*id, *volunteerID, *activityID, *date, *recordedBy, *hoursText, *notes)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(r.Output, "saved service record %s\n", stored.ID)
		return err
	case "list":
		return r.runList([]string{"records"})
	default:
		return fmt.Errorf("unknown record action %q", args[0])
	}
}

func (r *Runner) runReview(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("review action and record id are required")
	}
	approve := args[0] == "approve"
	if !approve && args[0] != "reject" {
		return fmt.Errorf("unknown review action %q", args[0])
	}
	reason := ""
	if len(args) > 2 {
		reason = strings.Join(args[2:], " ")
	}
	record, err := r.Ops.ReviewPending(args[1], "cli", approve, reason)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(r.Output, "record %s is %s\n", record.ID, record.Status)
	return err
}

func (r *Runner) runList(args []string) error {
	if len(args) == 0 || args[0] != "records" {
		return fmt.Errorf("list records is required")
	}
	flags := flag.NewFlagSet("list records", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	volunteerID := flags.String("volunteer", "", "volunteer")
	activityID := flags.String("activity", "", "activity")
	status := flags.String("status", "", "status")
	sortField := flags.String("sort", "date", "sort")
	descending := flags.Bool("desc", false, "descending")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	result, err := r.Ops.Search(query.Filter{VolunteerID: *volunteerID, ActivityID: *activityID, Status: *status}, domain.ParseSortSpec(*sortField, *descending), 0, 0)
	if err != nil {
		return err
	}
	return report.WriteRecords(r.Output, result)
}

func (r *Runner) runSummary(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("summary does not accept arguments")
	}
	result, _, err := r.Ops.Summary(query.Filter{})
	if err != nil {
		return err
	}
	volunteers, err := r.Intake.Catalog.ListVolunteers()
	if err != nil {
		return err
	}
	activities, err := r.Intake.Catalog.ListActivities()
	if err != nil {
		return err
	}
	return report.WriteSummary(r.Output, report.BuildSummary(result, volunteers, activities))
}
