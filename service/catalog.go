package service

import (
	"fmt"
	"strings"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

type CatalogService struct {
	store *storage.Store
}

func NewCatalogService(store *storage.Store) *CatalogService {
	return &CatalogService{store: store}
}

func (s *CatalogService) RegisterVolunteer(volunteer domain.Volunteer) (domain.Volunteer, error) {
	if err := domain.ValidateVolunteer(volunteer); err != nil {
		return domain.Volunteer{}, err
	}
	if _, err := s.store.GetVolunteer(volunteer.ID); err == nil {
		return domain.Volunteer{}, fmt.Errorf("volunteer %s already exists", volunteer.ID)
	}
	if err := s.store.SaveVolunteer(volunteer); err != nil {
		return domain.Volunteer{}, err
	}
	return volunteer, nil
}

func (s *CatalogService) UpdateVolunteer(volunteer domain.Volunteer) (domain.Volunteer, error) {
	if err := domain.ValidateVolunteer(volunteer); err != nil {
		return domain.Volunteer{}, err
	}
	if _, err := s.store.GetVolunteer(volunteer.ID); err != nil {
		return domain.Volunteer{}, fmt.Errorf("cannot update volunteer: %w", err)
	}
	if err := s.store.SaveVolunteer(volunteer); err != nil {
		return domain.Volunteer{}, err
	}
	return volunteer, nil
}

func (s *CatalogService) FindVolunteer(id string) (domain.Volunteer, error) {
	return s.store.GetVolunteer(strings.TrimSpace(id))
}

func (s *CatalogService) ListVolunteers() ([]domain.Volunteer, error) {
	return s.store.ListVolunteers()
}

func (s *CatalogService) CreateActivity(activity domain.Activity) (domain.Activity, error) {
	if err := domain.ValidateActivity(activity); err != nil {
		return domain.Activity{}, err
	}
	if _, err := s.store.GetActivity(activity.ID); err == nil {
		return domain.Activity{}, fmt.Errorf("activity %s already exists", activity.ID)
	}
	if err := s.store.SaveActivity(activity); err != nil {
		return domain.Activity{}, err
	}
	return activity, nil
}

func (s *CatalogService) UpdateActivity(activity domain.Activity) (domain.Activity, error) {
	if err := domain.ValidateActivity(activity); err != nil {
		return domain.Activity{}, err
	}
	if _, err := s.store.GetActivity(activity.ID); err != nil {
		return domain.Activity{}, fmt.Errorf("cannot update activity: %w", err)
	}
	if err := s.store.SaveActivity(activity); err != nil {
		return domain.Activity{}, err
	}
	return activity, nil
}

func (s *CatalogService) FindActivity(id string) (domain.Activity, error) {
	return s.store.GetActivity(strings.TrimSpace(id))
}

func (s *CatalogService) ListActivities() ([]domain.Activity, error) {
	return s.store.ListActivities()
}

func (s *CatalogService) OpenActivity(id string) (domain.Activity, error) {
	activity, err := s.FindActivity(id)
	if err != nil {
		return domain.Activity{}, err
	}
	if activity.Status == domain.ActivityClosed {
		return domain.Activity{}, fmt.Errorf("activity %s is closed", id)
	}
	activity.Status = domain.ActivityOpen
	return s.UpdateActivity(activity)
}

func (s *CatalogService) CloseActivity(id string) (domain.Activity, error) {
	activity, err := s.FindActivity(id)
	if err != nil {
		return domain.Activity{}, err
	}
	if activity.Status == domain.ActivityClosed {
		return activity, nil
	}
	activity.Status = domain.ActivityClosed
	return s.UpdateActivity(activity)
}

func (s *CatalogService) ActiveVolunteerCount() (int, error) {
	volunteers, err := s.ListVolunteers()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, volunteer := range volunteers {
		if volunteer.Active {
			count++
		}
	}
	return count, nil
}
