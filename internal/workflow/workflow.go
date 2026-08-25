package workflow

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/service"
	"example.com/noise-registry/internal/store"
	"fmt"
)

type Engine struct {
	Svc   *service.Service
	Store *store.Store
}

func New(s *store.Store) *Engine { return &Engine{Svc: service.New(s), Store: s} }
func (e *Engine) CreateReviewArchive(id, name string, l model.Level) (model.Record, error) {
	r, x := e.Svc.Register(id, name, l)
	if x != nil {
		return r, x
	}
	if r, x = e.Svc.Review(id); x != nil {
		return r, x
	}
	if r, x = e.Svc.Confirm(id); x != nil {
		return r, x
	}
	if r, x = e.Svc.Publish(id); x != nil {
		return r, x
	}
	return e.Svc.Archive(id)
}
func (e *Engine) SearchUpdatePublish(id string, l model.Level) (model.Record, error) {
	r, x := e.Svc.UpdateLevel(id, l)
	if x != nil {
		return r, x
	}
	if r.Status == model.StatusPublished {
		return r, nil
	}
	if r.Status == model.StatusApproved {
		return e.Svc.Publish(id)
	}
	return r, fmt.Errorf("record %s not ready", id)
}
func (e *Engine) ImportReport(rows []struct{ ID, Name, Level string }) (int, int, error) {
	for _, row := range rows {
		l, ok := model.ParseLevel(row.Level)
		if !ok {
			return 0, 1, fmt.Errorf("invalid level")
		}
		if _, e := e.Svc.Register(row.ID, row.Name, l); e != nil {
			return 0, 1, e
		}
	}
	return len(rows), 0, nil
}
func (e *Engine) StartWorkflow(id, owner string, records []string) error {
	return e.Store.PutWorkflow(model.Workflow{ID: id, Owner: owner, RecordIDs: records, Stage: "started", UpdatedAt: e.Store.Now()})
}
func (e *Engine) AdvanceWorkflow(id, stage string) error {
	w, err := e.Store.GetWorkflow(id)
	if err != nil {
		return err
	}
	w.Stage = stage
	w.UpdatedAt = e.Store.Now()
	return e.Store.PutWorkflow(w)
}
