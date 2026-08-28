package report

import (
	"fmt"
	"io"

	"volunteerhours/domain"
	"volunteerhours/query"
)

func WriteDailyTotals(w io.Writer, records []domain.ServiceRecord) error {
	if _, err := fmt.Fprintln(w, "Date\tRecords\tHours"); err != nil {
		return err
	}
	for _, total := range query.DailyTotals(records) {
		if _, err := fmt.Fprintf(w, "%s\t%d\t%.2f\n", total.Date, total.Count, total.Hours); err != nil {
			return err
		}
	}
	return nil
}

func WriteStatusTotals(w io.Writer, records []domain.ServiceRecord) error {
	if _, err := fmt.Fprintln(w, "Status\tRecords\tHours"); err != nil {
		return err
	}
	for _, total := range query.StatusTotals(records) {
		if _, err := fmt.Fprintf(w, "%s\t%d\t%.2f\n", total.Status, total.Count, total.Hours); err != nil {
			return err
		}
	}
	return nil
}

func RenderRecord(record domain.ServiceRecord, result query.Result) string {
	return fmt.Sprintf("%s | %s | %s | %s | %.2f | %s", record.ID, result.VolunteerName(record.VolunteerID), result.ActivityTitle(record.ActivityID), record.ServiceDate, record.Hours, record.Status)
}

func RenderActivity(activity domain.Activity) string {
	return fmt.Sprintf("%s: %s (%s-%s) %s", activity.ID, activity.Title, activity.StartDate, activity.EndDate, activity.Status)
}

func RenderVolunteer(volunteer domain.Volunteer) string {
	state := "inactive"
	if volunteer.Active {
		state = "active"
	}
	return fmt.Sprintf("%s: %s <%s> %s", volunteer.ID, volunteer.DisplayName(), volunteer.Email, state)
}
