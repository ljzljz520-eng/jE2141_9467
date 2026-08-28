package workflow

import (
	"fmt"

	"volunteerhours/domain"
	"volunteerhours/storage"
)

type Recovery struct {
	Path string
}

func NewRecovery(path string) Recovery {
	return Recovery{Path: path}
}

func (r Recovery) ReopenAndInspect() (*storage.Store, int, error) {
	store, err := storage.Open(r.Path)
	if err != nil {
		return nil, 0, err
	}
	count, err := store.RecordCount()
	if err != nil {
		store.Close()
		return nil, 0, err
	}
	return store, count, nil
}

func (r Recovery) VerifyEntities(store *storage.Store, ids []string) error {
	for _, id := range ids {
		if _, err := store.GetServiceRecord(id); err != nil {
			return fmt.Errorf("recovery missing ServiceRecord %s: %w", id, err)
		}
	}
	return nil
}

func (r Recovery) ExportSnapshot(store *storage.Store) ([]domain.ServiceRecord, error) {
	records, err := store.ListServiceRecords()
	if err != nil {
		return nil, err
	}
	return append([]domain.ServiceRecord(nil), records...), nil
}

func (r Recovery) ReopenWithValidation() (*storage.Store, error) {
	store, _, err := r.ReopenAndInspect()
	if err != nil {
		return nil, err
	}
	if _, err := store.ListVolunteers(); err != nil {
		store.Close()
		return nil, err
	}
	if _, err := store.ListActivities(); err != nil {
		store.Close()
		return nil, err
	}
	return store, nil
}

func (r Recovery) Close(store *storage.Store) error {
	if store == nil {
		return nil
	}
	return store.Close()
}
