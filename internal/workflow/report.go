package workflow

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/store"
	"strings"
)

type Report struct {
	Total, Published, Archived int
	Names                      []string
}

func BuildReport(s *store.Store) (Report, error) {
	rs, e := s.ListRecords("")
	if e != nil {
		return Report{}, e
	}
	r := Report{Total: len(rs)}
	for _, x := range rs {
		if x.Status == model.StatusPublished {
			r.Published++
		}
		if x.Status == model.StatusArchived {
			r.Archived++
		}
		r.Names = append(r.Names, x.Name)
	}
	return r, nil
}
func RenderReport(r Report) string {
	return strings.Join([]string{"total=" + itoa(r.Total), "published=" + itoa(r.Published), "archived=" + itoa(r.Archived), "names=" + strings.Join(r.Names, "|")}, " ")
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(rune('0'+n%10)) + out
		n /= 10
	}
	return out
}
func EmptyReport() Report { return Report{Names: []string{}} }
func MergeReports(a, b Report) Report {
	a.Total += b.Total
	a.Published += b.Published
	a.Archived += b.Archived
	a.Names = append(a.Names, b.Names...)
	return a
}
