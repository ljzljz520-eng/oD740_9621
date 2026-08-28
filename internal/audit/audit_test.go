package audit

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"testing"
)

func TestAuditHistory(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	s.PutAudit(model.NewAudit("a", "r", "registered", "x", "d", s.Now()))
	r := New(s)
	events, e := r.History("r")
	if e != nil || len(events) != 1 {
		t.Fatal(e)
	}
	if !ContainsAction(events, "registered") {
		t.Fatal("action")
	}
}
