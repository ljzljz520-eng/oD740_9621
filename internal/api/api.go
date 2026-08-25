package api

import (
	"encoding/json"
	"example.com/noise-registry/internal/service"
	"net/http"
	"strings"
)

type Server struct{ Svc *service.Service }

func New(s *service.Service) *Server { return &Server{Svc: s} }
func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/health", s.health)
	m.HandleFunc("/records", s.records)
	return m
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	rs, e := s.Svc.Search(strings.TrimSpace(q))
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(rs)
}
func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(v)
}
