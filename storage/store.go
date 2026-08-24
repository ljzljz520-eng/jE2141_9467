package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"
)

var (
	volunteerBucket = []byte("volunteers")
	activityBucket  = []byte("activities")
	recordBucket    = []byte("service_records")
	statusBucket    = []byte("service_statuses")
	auditBucket     = []byte("audit_entries")
)

type Store struct {
	path string
	db   *bbolt.DB
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt database: %w", err)
	}
	store := &Store{path: path, db: db}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{volunteerBucket, activityBucket, recordBucket, statusBucket, auditBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) healthy() error {
	if s.db == nil {
		return errors.New("storage is closed")
	}
	return nil
}

func (s *Store) write(bucket []byte, key string, value []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return err
	}
	if key == "" {
		return errors.New("storage key is required")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), value)
	})
}

func (s *Store) read(bucket []byte, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return nil, err
	}
	var copyValue []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return fmt.Errorf("record %q not found", key)
		}
		copyValue = append([]byte(nil), value...)
		return nil
	})
	return copyValue, err
}

func (s *Store) remove(bucket []byte, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Delete([]byte(key))
	})
}

func (s *Store) list(bucket []byte) (map[string][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return nil, err
	}
	items := make(map[string][]byte)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(k, v []byte) error {
			items[string(k)] = append([]byte(nil), v...)
			return nil
		})
	})
	return items, err
}

func (s *Store) Count(bucket []byte) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return 0, err
	}
	count := 0
	err := s.db.View(func(tx *bbolt.Tx) error {
		count = tx.Bucket(bucket).Stats().KeyN
		return nil
	})
	return count, err
}

func (s *Store) DeleteAll() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.healthy(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{volunteerBucket, activityBucket, recordBucket, statusBucket, auditBucket} {
			if err := tx.Bucket(bucket).ForEach(func(k, _ []byte) error { return tx.Bucket(bucket).Delete(k) }); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) SaveEntity(entity string, id string, value []byte) error {
	switch entity {
	case "Volunteer":
		return s.write(volunteerBucket, id, value)
	case "Activity":
		return s.write(activityBucket, id, value)
	case "ServiceRecord":
		return s.write(recordBucket, id, value)
	case "ServiceStatus":
		return s.write(statusBucket, id, value)
	case "AuditEntry":
		return s.write(auditBucket, id, value)
	default:
		return fmt.Errorf("unsupported entity %s", entity)
	}
}
