package api

import (
	"net/http"
	"strings"
)

func RouteName(path string) string {
	p := strings.Trim(path, "/")
	if p == "" {
		return "root"
	}
	parts := strings.Split(p, "/")
	return parts[0]
}
func IsRecordRoute(path string) bool    { return RouteName(path) == "records" }
func IsHealthRoute(path string) bool    { return RouteName(path) == "health" }
func AllowOrigin(w http.ResponseWriter) { w.Header().Set("access-control-allow-origin", "*") }
func RequestID(r *http.Request) string {
	if v := r.Header.Get("x-request-id"); v != "" {
		return v
	}
	return "local-request"
}
