package flow033

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/service"
	"example.com/noise-registry/internal/store"
	"testing"
)

func Test740BusinessRegression(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	x := service.NewProcessor(service.New(s))
	x.Svc.Register("one", "First", model.LevelLight)
	x.Svc.Register("two", "Second", model.LevelStrong)
	x.Svc.Review("one")
	x.Svc.Confirm("one")
	x.Svc.Review("two")
	x.Svc.Confirm("two")
	x.Prepare("one", model.LevelLight)
	x.Prepare("two", model.LevelStrong)
	x.Publish("one")
	x.Publish("two")
	a, _ := s.GetRecord("one")
	b, _ := s.GetRecord("two")
	if a.Level == b.Level {
		t.Fatalf("levels merged: %s", a.Level)
	}
}
