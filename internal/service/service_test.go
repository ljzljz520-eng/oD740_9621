package service

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"testing"
)

func TestLifecycleRules(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	x := New(s)
	if _, e = x.Register("r", "N", model.LevelStrong); e != nil {
		t.Fatal(e)
	}
	if _, e = x.Publish("r"); e == nil {
		t.Fatal("draft published")
	}
	if _, e = x.Review("r"); e != nil {
		t.Fatal(e)
	}
}
