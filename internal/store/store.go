package store

import (
	"encoding/json"
	"errors"
	"example.com/noise-registry/internal/model"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var recordsBucket = []byte("records")
var auditsBucket = []byte("audits")
var workflowsBucket = []byte("workflows")
var attachmentsBucket = []byte("attachments")

type Store struct {
	db  *bbolt.DB
	mu  sync.RWMutex
	now func() time.Time
}

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(path, 0600, nil)
	if e != nil {
		return nil, e
	}
	s := &Store{db: db, now: func() time.Time { return time.Unix(1700000000, 0) }}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, n := range [][]byte{recordsBucket, auditsBucket, workflowsBucket, attachmentsBucket} {
			if _, x := tx.CreateBucketIfNotExists(n); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func OpenTemp() (*Store, string, error) {
	f, e := os.CreateTemp("", "noise-registry-*.db")
	if e != nil {
		return nil, "", e
	}
	p := f.Name()
	f.Close()
	s, e := Open(p)
	return s, p, e
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func encode(v any) ([]byte, error) { return json.Marshal(v) }
func decode(b []byte, v any) error { return json.Unmarshal(b, v) }
func (s *Store) PutRecord(r model.Record) error {
	if e := r.Validate(); e != nil {
		return e
	}
	b, e := encode(r)
	if e != nil {
		return e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(recordsBucket).Put([]byte(r.ID), b) })
}
func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(recordsBucket).Get([]byte(id))
		if v == nil {
			return errors.New("record not found")
		}
		return decode(v, &r)
	})
	return r, e
}
func (s *Store) ListRecords(q string) ([]model.Record, error) {
	var out []model.Record
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordsBucket).ForEach(func(_, v []byte) error {
			var r model.Record
			if e := decode(v, &r); e != nil {
				return e
			}
			if q == "" || strings.Contains(strings.ToLower(r.Name), strings.ToLower(q)) {
				out = append(out, r)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, e
}
func (s *Store) PutAudit(a model.AuditEvent) error {
	b, e := encode(a)
	if e != nil {
		return e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(auditsBucket).Put([]byte(a.ID), b) })
}
func (s *Store) Audits(recordID string) ([]model.AuditEvent, error) {
	var out []model.AuditEvent
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditsBucket).ForEach(func(_, v []byte) error {
			var a model.AuditEvent
			if e := decode(v, &a); e != nil {
				return e
			}
			if recordID == "" || a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out, e
}
func (s *Store) PutWorkflow(w model.Workflow) error {
	if e := w.Valid(); e != nil {
		return e
	}
	b, e := encode(w)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(workflowsBucket).Put([]byte(w.ID), b) })
}
func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var w model.Workflow
	e := s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(workflowsBucket).Get([]byte(id))
		if v == nil {
			return fmt.Errorf("workflow %s not found", id)
		}
		return decode(v, &w)
	})
	return w, e
}
func (s *Store) PutAttachment(a model.Attachment) error {
	if e := a.Valid(); e != nil {
		return e
	}
	b, e := encode(a)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(attachmentsBucket).Put([]byte(a.ID), b) })
}
func (s *Store) Attachments(recordID string) ([]model.Attachment, error) {
	var out []model.Attachment
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(attachmentsBucket).ForEach(func(_, v []byte) error {
			var a model.Attachment
			if e := decode(v, &a); e != nil {
				return e
			}
			if recordID == "" || a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	return out, e
}
func (s *Store) Now() time.Time { return s.now() }
