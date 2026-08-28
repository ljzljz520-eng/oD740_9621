package audit

import (
	"example.com/noise-registry/internal/model"
	"sort"
	"strings"
)

func Actions(events []model.AuditEvent) []string {
	out := []string{}
	for _, e := range events {
		out = append(out, e.Action)
	}
	return out
}
func ByActor(events []model.AuditEvent, actor string) []model.AuditEvent {
	out := []model.AuditEvent{}
	for _, e := range events {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}
func ContainsAction(events []model.AuditEvent, action string) bool {
	for _, e := range events {
		if strings.EqualFold(e.Action, action) {
			return true
		}
	}
	return false
}
func Chronological(events []model.AuditEvent) []model.AuditEvent {
	sort.SliceStable(events, func(i, j int) bool { return events[i].At.Before(events[j].At) })
	return events
}
func CountAction(events []model.AuditEvent, action string) int {
	n := 0
	for _, e := range events {
		if e.Action == action {
			n++
		}
	}
	return n
}
