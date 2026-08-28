package importer

import (
	"example.com/noise-registry/internal/model"
	"example.com/noise-registry/internal/service"
	"fmt"
	"strings"
)

type Row struct{ ID, Name, Level string }
type Result struct {
	Accepted, Rejected int
	Errors             []string
	Records            []model.Record
}

func Batch(s *service.Service, rows []Row) Result {
	out := Result{}
	for _, row := range rows {
		l, ok := model.ParseLevel(strings.TrimSpace(row.Level))
		if row.ID == "" || row.Name == "" || !ok {
			out.Rejected++
			out.Errors = append(out.Errors, fmt.Sprintf("%s rejected", row.ID))
			continue
		}
		r, e := s.Register(row.ID, row.Name, l)
		if e != nil {
			out.Rejected++
			out.Errors = append(out.Errors, e.Error())
			continue
		}
		out.Accepted++
		out.Records = append(out.Records, r)
	}
	return out
}
func ValidateRows(rows []Row) []string {
	errs := []string{}
	for _, r := range rows {
		if r.ID == "" {
			errs = append(errs, "id required")
		}
		if r.Name == "" {
			errs = append(errs, "name required")
		}
	}
	return errs
}
func CSV(rows []Row) string {
	out := "id,name,level\n"
	for _, r := range rows {
		out += fmt.Sprintf("%s,%s,%s\n", r.ID, r.Name, r.Level)
	}
	return out
}
