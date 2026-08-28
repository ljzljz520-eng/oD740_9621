package store

import (
	"example.com/noise-registry/internal/model"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	s, p, e := OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	r := model.NewRecord("persist", "Persisted", model.LevelLight, s.Now())
	if e = s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetRecord("persist")
	if e != nil || got.Name != "Persisted" {
		t.Fatalf("%v %#v", e, got)
	}
}
