package service

import (
	"errors"
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"fmt"
)

type Service struct {
	Store *store.Store
	Actor string
}

func New(s *store.Store) *Service { return &Service{Store: s, Actor: "operator"} }
func (s *Service) Register(id, name string, l model.Level) (model.Record, error) {
	r := model.NewRecord(id, name, l, s.Store.Now())
	if e := s.Store.PutRecord(r); e != nil {
		return r, e
	}
	s.audit(r, "registered", "draft created")
	return r, nil
}
func (s *Service) audit(r model.Record, action, detail string) error {
	a := model.NewAudit(fmt.Sprintf("%s-%s-%d", r.ID, action, r.Version), r.ID, action, s.Actor, detail, s.Store.Now())
	return s.Store.PutAudit(a)
}
func (s *Service) Review(id string) (model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return r, e
	}
	if r.Status != model.StatusDraft {
		return r, errors.New("only draft can review")
	}
	r.Advance(model.StatusReviewed, s.Store.Now())
	e = s.Store.PutRecord(r)
	if e == nil {
		e = s.audit(r, "reviewed", "metadata checked")
	}
	return r, e
}
func (s *Service) Confirm(id string) (model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return r, e
	}
	if r.Status != model.StatusReviewed {
		return r, errors.New("review required")
	}
	r.Advance(model.StatusApproved, s.Store.Now())
	e = s.Store.PutRecord(r)
	if e == nil {
		e = s.audit(r, "approved", "review confirmed")
	}
	return r, e
}
func (s *Service) Publish(id string) (model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return r, e
	}
	if !r.Status.CanPublish() {
		return r, errors.New("approval required")
	}
	r.Advance(model.StatusPublished, s.Store.Now())
	e = s.Store.PutRecord(r)
	if e == nil {
		e = s.audit(r, "published", "profile published")
	}
	return r, e
}
func (s *Service) Archive(id string) (model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return r, e
	}
	if r.Status != model.StatusPublished {
		return r, errors.New("only published can archive")
	}
	r.Advance(model.StatusArchived, s.Store.Now())
	e = s.Store.PutRecord(r)
	if e == nil {
		e = s.audit(r, "archived", "profile archived")
	}
	return r, e
}
func (s *Service) Search(q string) ([]model.Record, error) { return s.Store.ListRecords(q) }
func (s *Service) UpdateLevel(id string, l model.Level) (model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return r, e
	}
	if !l.Valid() {
		return r, errors.New("invalid level")
	}
	if r.Status.Terminal() {
		return r, errors.New("archived")
	}
	r.Level = l
	r.Version++
	r.Touch(s.Store.Now())
	e = s.Store.PutRecord(r)
	if e == nil {
		e = s.audit(r, "updated", "level changed")
	}
	return r, e
}
func (s *Service) Attach(id, name, checksum string, size int64) (model.Attachment, error) {
	a := model.Attachment{ID: fmt.Sprintf("%s-att-%d", id, size), RecordID: id, Name: name, Checksum: checksum, Size: size, CreatedAt: s.Store.Now()}
	e := s.Store.PutAttachment(a)
	return a, e
}
