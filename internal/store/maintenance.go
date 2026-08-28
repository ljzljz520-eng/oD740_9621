package store

import (
	"errors"
	"example.com/noise-registry/internal/model"
	"go.etcd.io/bbolt"
	"time"
)

func (s *Store) DeleteRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(recordsBucket).Delete([]byte(id)) })
}
func (s *Store) Exists(id string) bool { _, e := s.GetRecord(id); return e == nil }
func (s *Store) ReplaceRecord(r model.Record) error {
	if !s.Exists(r.ID) {
		return errors.New("record missing")
	}
	return s.PutRecord(r)
}
func (s *Store) RecordsByTime(after, before time.Time) ([]model.Record, error) {
	rs, e := s.ListRecords("")
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, r := range rs {
		if !r.CreatedAt.Before(after) && r.CreatedAt.Before(before) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) Count() int {
	rs, e := s.ListRecords("")
	if e != nil {
		return 0
	}
	return len(rs)
}
func (s *Store) AuditCount() int {
	as, e := s.Audits("")
	if e != nil {
		return 0
	}
	return len(as)
}
func (s *Store) Health() error {
	if s == nil || s.db == nil {
		return errors.New("closed")
	}
	return s.db.View(func(*bbolt.Tx) error { return nil })
}
