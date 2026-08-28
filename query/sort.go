package query

import (
	"sort"
	"strings"

	"volunteerhours/domain"
)

type Group struct {
	Key       string
	Label     string
	Records   []domain.ServiceRecord
	TotalHour float64
}

func GroupByVolunteer(result Result) []Group {
	groups := make(map[string]*Group)
	for _, record := range result.Records {
		group, ok := groups[record.VolunteerID]
		if !ok {
			group = &Group{Key: record.VolunteerID, Label: result.VolunteerName(record.VolunteerID)}
			groups[record.VolunteerID] = group
		}
		group.Records = append(group.Records, record)
		if record.Status != domain.RecordRejected {
			group.TotalHour += record.Hours
		}
	}
	return orderedGroups(groups)
}

func GroupByActivity(result Result) []Group {
	groups := make(map[string]*Group)
	for _, record := range result.Records {
		group, ok := groups[record.ActivityID]
		if !ok {
			group = &Group{Key: record.ActivityID, Label: result.ActivityTitle(record.ActivityID)}
			groups[record.ActivityID] = group
		}
		group.Records = append(group.Records, record)
		if record.Status != domain.RecordRejected {
			group.TotalHour += record.Hours
		}
	}
	return orderedGroups(groups)
}

func orderedGroups(groups map[string]*Group) []Group {
	result := make([]Group, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(i, j int) bool {
		left, right := strings.ToLower(result[i].Label), strings.ToLower(result[j].Label)
		if left == right {
			return result[i].Key < result[j].Key
		}
		return left < right
	})
	return result
}

func Limit(result Result, offset, limit int) Result {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if offset >= len(result.Records) {
		result.Records = nil
		return result
	}
	end := len(result.Records)
	if limit > 0 && offset+limit < end {
		end = offset + limit
	}
	result.Records = result.Records[offset:end]
	return result
}
