package service

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"os"
	"testing"
)

func TestProcessorLevelsNotCollapsed(t *testing.T) {
	s, path, err := store.OpenTemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	defer s.Close()
	svc := New(s)

	// two records, two distinct levels that must not override each other
	if _, err := svc.Register("r1", "alpha", model.LevelLight); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Register("r2", "beta", model.LevelStrong); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Review("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Confirm("r1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Review("r2"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Confirm("r2"); err != nil {
		t.Fatal(err)
	}

	p := NewProcessor(svc)
	if err := p.Prepare("r1", model.LevelLight); err != nil {
		t.Fatal(err)
	}
	if err := p.Prepare("r2", model.LevelStrong); err != nil {
		t.Fatal(err)
	}

	// Publish the last-prepared record first: must not leak its level onto r1.
	if r, err := p.Publish("r2"); err != nil {
		t.Fatal(err)
	} else if r.Level != model.LevelStrong {
		t.Fatalf("r2 level = %s, want strong", r.Level)
	}
	// r1 must still be light — the bug collapsed it into the shared lastLevel (strong).
	if r, err := p.Publish("r1"); err != nil {
		t.Fatal(err)
	} else if r.Level != model.LevelLight {
		t.Fatalf("r1 level = %s, want light (levels were collapsed)", r.Level)
	}

	// persisted values
	if r, err := s.GetRecord("r1"); err != nil {
		t.Fatal(err)
	} else if r.Level != model.LevelLight {
		t.Fatalf("persisted r1 level = %s, want light", r.Level)
	}
	if r, err := s.GetRecord("r2"); err != nil {
		t.Fatal(err)
	} else if r.Level != model.LevelStrong {
		t.Fatalf("persisted r2 level = %s, want strong", r.Level)
	}
}
