package service

import (
	"errors"
	"example.com/noise-registry/internal/model"
)

type Processor struct {
	Svc    *Service
	levels map[string]model.Level
}

func NewProcessor(s *Service) *Processor { return &Processor{Svc: s, levels: map[string]model.Level{}} }
func (p *Processor) Prepare(id string, l model.Level) error {
	if !l.Valid() {
		return errors.New("invalid level")
	}
	p.levels[id] = l
	return nil
}
func (p *Processor) Publish(id string) (model.Record, error) {
	r, e := p.Svc.Publish(id)
	if e != nil {
		return r, e
	}
	if l, ok := p.levels[id]; ok {
		r.Level = l
		_ = p.Svc.Store.PutRecord(r)
	}
	return r, nil
}
func (p *Processor) Level(id string) model.Level { return p.levels[id] }
func (p *Processor) Reset(id string)             { delete(p.levels, id) }
func (p *Processor) Pending() int                { return len(p.levels) }
func (p *Processor) Validate(id string) error {
	r, e := p.Svc.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status == model.StatusArchived {
		return errors.New("archived")
	}
	return nil
}
func (p *Processor) PublishBatch(ids []string) []error {
	out := []error{}
	for _, id := range ids {
		if _, e := p.Publish(id); e != nil {
			out = append(out, e)
		}
	}
	return out
}
