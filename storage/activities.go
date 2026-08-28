package storage

import (
	"fmt"
	"sort"

	"volunteerhours/domain"
)

func (s *Store) SaveActivity(activity domain.Activity) error {
	data, err := domain.Encode(activity)
	if err != nil {
		return fmt.Errorf("encode Activity: %w", err)
	}
	return s.write(activityBucket, activity.ID, data)
}

func (s *Store) GetActivity(id string) (domain.Activity, error) {
	data, err := s.read(activityBucket, id)
	if err != nil {
		return domain.Activity{}, err
	}
	var activity domain.Activity
	if err := domain.Decode(data, &activity); err != nil {
		return domain.Activity{}, fmt.Errorf("decode Activity: %w", err)
	}
	return activity, nil
}

func (s *Store) ListActivities() ([]domain.Activity, error) {
	items, err := s.list(activityBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Activity, 0, len(items))
	for _, data := range items {
		var activity domain.Activity
		if err := domain.Decode(data, &activity); err != nil {
			return nil, fmt.Errorf("decode Activity list: %w", err)
		}
		result = append(result, activity)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StartDate < result[j].StartDate })
	return result, nil
}

func (s *Store) DeleteActivity(id string) error {
	return s.remove(activityBucket, id)
}

func (s *Store) ActivityCount() (int, error) {
	return s.Count(activityBucket)
}
