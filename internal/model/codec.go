package model

import (
	"encoding/json"
	"sort"
	"strings"
)

func EncodeRecord(r Record) ([]byte, error)    { return json.Marshal(r) }
func DecodeRecord(b []byte) (Record, error)    { var r Record; e := json.Unmarshal(b, &r); return r, e }
func EncodeAudit(a AuditEvent) ([]byte, error) { return json.Marshal(a) }
func DecodeAudit(b []byte) (AuditEvent, error) {
	var a AuditEvent
	e := json.Unmarshal(b, &a)
	return a, e
}
func NormalizeStatus(v string) Status {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case string(StatusDraft):
		return StatusDraft
	case string(StatusReviewed):
		return StatusReviewed
	case string(StatusApproved):
		return StatusApproved
	case string(StatusPublished):
		return StatusPublished
	case string(StatusArchived):
		return StatusArchived
	}
	return StatusDraft
}
func SortRecords(rs []Record) []Record {
	sort.SliceStable(rs, func(i, j int) bool {
		if rs[i].Name == rs[j].Name {
			return rs[i].ID < rs[j].ID
		}
		return rs[i].Name < rs[j].Name
	})
	return rs
}
func CloneRecords(rs []Record) []Record { out := make([]Record, len(rs)); copy(out, rs); return out }
func IDs(rs []Record) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.ID)
	}
	return out
}
func HasLevel(rs []Record, l Level) bool {
	for _, r := range rs {
		if r.Level == l {
			return true
		}
	}
	return false
}
func GroupByStatus(rs []Record) map[Status][]Record {
	out := map[Status][]Record{}
	for _, r := range rs {
		out[r.Status] = append(out[r.Status], r)
	}
	return out
}
