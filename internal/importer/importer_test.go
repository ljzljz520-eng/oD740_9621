package importer

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/service"
	"example.com/noise-registry/internal/store"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	s, _, e := store.OpenTemp()
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := Batch(service.New(s), []Row{{"a", "A", "light"}, {"b", "", "strong"}})
	if r.Accepted != 1 || r.Rejected != 1 {
		t.Fatalf("%+v", r)
	}
}
func TestDeterministicChecksum(t *testing.T) {
	if Checksum(Row{"a", "A", "light"}) != Checksum(Row{"a", "A", "light"}) {
		t.Fatal("unstable")
	}
	if !Rejected(Result{Rejected: 1}) {
		t.Fatal("expected rejected")
	}
	_ = model.LevelLight
}
