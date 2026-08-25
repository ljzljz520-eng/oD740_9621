package workflow

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r, e := New(s).CreateReviewArchive("r", "Clean speech", model.LevelBalanced)
	if e != nil || r.Status != model.StatusArchived {
		t.Fatalf("%v %v", e, r.Status)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	x := New(s)
	r, _ := x.Svc.Register("r", "Clean", model.LevelLight)
	x.Svc.Review(r.ID)
	x.Svc.Confirm(r.ID)
	if _, e = x.SearchUpdatePublish(r.ID, model.LevelStrong); e != nil {
		t.Fatal(e)
	}
}
