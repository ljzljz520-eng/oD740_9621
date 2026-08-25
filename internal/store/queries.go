package store

import (
	"example.com/noise-registry/internal/model"
	"strings"
)

func FilterStatus(rs []model.Record, status model.Status) []model.Record {
	out := make([]model.Record, 0)
	for _, r := range rs {
		if status == "" || r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
func MatchName(r model.Record, q string) bool {
	return strings.Contains(strings.ToLower(r.Name), strings.ToLower(q))
}
func CountPublished(rs []model.Record) int {
	n := 0
	for _, r := range rs {
		if r.Status == model.StatusPublished {
			n++
		}
	}
	return n
}
