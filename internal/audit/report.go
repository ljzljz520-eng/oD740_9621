package audit

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"fmt"
	"strings"
)

type Reporter struct{ Store *store.Store }

func New(s *store.Store) *Reporter                                { return &Reporter{Store: s} }
func (r *Reporter) History(id string) ([]model.AuditEvent, error) { return r.Store.Audits(id) }
func (r *Reporter) Summary(id string) (string, error) {
	events, e := r.History(id)
	if e != nil {
		return "", e
	}
	parts := make([]string, 0, len(events))
	for _, a := range events {
		parts = append(parts, a.Action)
	}
	return fmt.Sprintf("%s: %s", id, strings.Join(parts, ",")), nil
}
func (r *Reporter) Latest(id string) (model.AuditEvent, error) {
	events, e := r.History(id)
	if e != nil || len(events) == 0 {
		return model.AuditEvent{}, fmt.Errorf("no audit for %s", id)
	}
	return events[len(events)-1], nil
}
