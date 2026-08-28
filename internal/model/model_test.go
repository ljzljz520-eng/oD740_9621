package model

import (
	"testing"
	"time"
)

func TestRecordValidation(t *testing.T) {
	r := NewRecord("r1", "Voice", LevelBalanced, time.Unix(1, 0))
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
	if Stronger(LevelLight, LevelStrong) != LevelStrong {
		t.Fatal("level")
	}
}
