package main

import (
	"example.com/noise-registry/internal/store"
	"testing"
)

func TestLocalStartupDependencies(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Health(); e != nil {
		t.Fatal(e)
	}
	s.Close()
}
