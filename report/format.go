package report

import (
	"fmt"
	"io"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/query"
)

func WriteSummary(w io.Writer, summary Summary) error {
	if _, err := fmt.Fprintf(w, "Volunteer hours summary\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Records: %d  Volunteers: %d  Activities: %d\n", summary.RecordCount, summary.VolunteerCount, summary.ActivityCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Total hours: %.2f\n", summary.TotalHours); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Status: %s\n", summary.StatusLabel()); err != nil {
		return err
	}
	if top, ok := summary.TopVolunteer(); ok {
		if _, err := fmt.Fprintf(w, "Top volunteer: %s (%.2f hours)\n", top.Label, top.TotalHour); err != nil {
			return err
		}
	}
	return nil
}

func WriteRecords(w io.Writer, result query.Result) error {
	if len(result.Records) == 0 {
		_, err := fmt.Fprintln(w, "No service records found")
		return err
	}
	if _, err := fmt.Fprintln(w, "ID\tVolunteer\tActivity\tDate\tHours\tStatus\tRecorded by"); err != nil {
		return err
	}
	for _, record := range result.Records {
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.2f\t%s\t%s\n", record.ID, result.VolunteerName(record.VolunteerID), result.ActivityTitle(record.ActivityID), record.ServiceDate, record.Hours, record.Status, record.RecordedBy); err != nil {
			return err
		}
	}
	return nil
}

func WriteVolunteers(w io.Writer, volunteers []domain.Volunteer) error {
	for _, volunteer := range volunteers {
		active := "inactive"
		if volunteer.Active {
			active = "active"
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", volunteer.ID, volunteer.DisplayName(), volunteer.Email, active); err != nil {
			return err
		}
	}
	return nil
}

func WriteActivities(w io.Writer, activities []domain.Activity) error {
	for _, activity := range activities {
		line := strings.Join([]string{activity.ID, activity.Title, activity.Location, activity.StartDate, activity.EndDate, activity.Status}, "\t")
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
