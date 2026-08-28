package storage

import (
	"fmt"
	"sort"

	"volunteerhours/domain"
)

func (s *Store) SaveVolunteer(volunteer domain.Volunteer) error {
	data, err := domain.Encode(volunteer)
	if err != nil {
		return fmt.Errorf("encode Volunteer: %w", err)
	}
	return s.write(volunteerBucket, volunteer.ID, data)
}

func (s *Store) GetVolunteer(id string) (domain.Volunteer, error) {
	data, err := s.read(volunteerBucket, id)
	if err != nil {
		return domain.Volunteer{}, err
	}
	var volunteer domain.Volunteer
	if err := domain.Decode(data, &volunteer); err != nil {
		return domain.Volunteer{}, fmt.Errorf("decode Volunteer: %w", err)
	}
	return volunteer, nil
}

func (s *Store) ListVolunteers() ([]domain.Volunteer, error) {
	items, err := s.list(volunteerBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Volunteer, 0, len(items))
	for _, data := range items {
		var volunteer domain.Volunteer
		if err := domain.Decode(data, &volunteer); err != nil {
			return nil, fmt.Errorf("decode Volunteer list: %w", err)
		}
		result = append(result, volunteer)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *Store) DeleteVolunteer(id string) error {
	return s.remove(volunteerBucket, id)
}

func (s *Store) VolunteerCount() (int, error) {
	return s.Count(volunteerBucket)
}
