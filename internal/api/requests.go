package api

import (
	"encoding/json"
	"example.com/noise-registry/internal/model"
	"net/http"
)

type CreateRequest struct {
	ID, Name string
	Level    model.Level
}

func DecodeCreate(r *http.Request) (CreateRequest, error) {
	var q CreateRequest
	e := json.NewDecoder(r.Body).Decode(&q)
	return q, e
}
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
func WriteError(w http.ResponseWriter, err error) { http.Error(w, err.Error(), StatusCode(err)) }
func MethodAllowed(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodPost
}
