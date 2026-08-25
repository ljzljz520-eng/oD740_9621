package service

import (
	"errors"
	"example.com/noise-registry/internal/model"
)

func AllowedTransition(from, to model.Status) bool {
	switch from {
	case model.StatusDraft:
		return to == model.StatusReviewed
	case model.StatusReviewed:
		return to == model.StatusApproved
	case model.StatusApproved:
		return to == model.StatusPublished
	case model.StatusPublished:
		return to == model.StatusArchived
	}
	return false
}
func RequirePublished(r model.Record) error {
	if r.Status != model.StatusPublished {
		return errors.New("published record required")
	}
	return nil
}
func NormalizeName(v string) string {
	if v == "" {
		return ""
	}
	return v
}
