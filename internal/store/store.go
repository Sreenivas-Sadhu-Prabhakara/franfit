// Package store is the in-memory lead store with JSON snapshot persistence.
// Every write saves the full state to ./data/store.json; the file is loaded
// on boot if present.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Lead is one captured "request intro" submission.
type Lead struct {
	ID             string    `json:"id"`
	BrandID        string    `json:"brandId"`
	BrandName      string    `json:"brandName"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	BudgetL        float64   `json:"budgetL"`
	FitScore       int       `json:"fitScore"`
	City           string    `json:"city"`
	NotificationID string    `json:"notificationId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type snapshot struct {
	Seq   int    `json:"seq"`
	Leads []Lead `json:"leads"`
}

// Store guards all mutable state behind one mutex.
type Store struct {
	mu   sync.Mutex
	path string
	data snapshot
}

// Open loads the snapshot at path if it exists, or starts empty.
func Open(path string) (*Store, error) {
	s := &Store{path: path}
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read snapshot: %w", err)
	}
	if err := json.Unmarshal(raw, &s.data); err != nil {
		return nil, fmt.Errorf("parse snapshot %s: %w", path, err)
	}
	return s, nil
}

// AddLead assigns an id and timestamp, appends the lead and persists.
func (s *Store) AddLead(l Lead) (Lead, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.Seq++
	l.ID = fmt.Sprintf("lead-%04d", s.data.Seq)
	l.CreatedAt = time.Now().UTC()
	s.data.Leads = append(s.data.Leads, l)
	if err := s.saveLocked(); err != nil {
		return Lead{}, err
	}
	return l, nil
}

// SetNotificationID records the provider message id for a lead and persists.
func (s *Store) SetNotificationID(leadID, notifID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.data.Leads {
		if s.data.Leads[i].ID == leadID {
			s.data.Leads[i].NotificationID = notifID
			return s.saveLocked()
		}
	}
	return fmt.Errorf("lead %s not found", leadID)
}

// Leads returns a copy of all leads, newest first.
func (s *Store) Leads() []Lead {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Lead, len(s.data.Leads))
	for i, l := range s.data.Leads {
		out[len(s.data.Leads)-1-i] = l
	}
	return out
}

func (s *Store) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
