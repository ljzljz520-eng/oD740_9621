package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Level string

const (
	LevelLight    Level = "light"
	LevelBalanced Level = "balanced"
	LevelStrong   Level = "strong"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusReviewed  Status = "reviewed"
	StatusApproved  Status = "approved"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Record struct {
	ID, Name             string
	Level                Level
	Status               Status
	CreatedAt, UpdatedAt time.Time
	Version              int
	Notes                string
}
type AuditEvent struct {
	ID, RecordID, Action, Actor, Detail string
	At                                  time.Time
}
type Workflow struct {
	ID           string
	RecordIDs    []string
	Stage, Owner string
	UpdatedAt    time.Time
}
type Attachment struct {
	ID, RecordID, Name, Checksum string
	Size                         int64
	CreatedAt                    time.Time
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("id required")
	}
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name required")
	}
	if !r.Level.Valid() {
		return fmt.Errorf("invalid level %q", r.Level)
	}
	if r.Version < 1 {
		return errors.New("version required")
	}
	return nil
}
func (l Level) Valid() bool {
	switch l {
	case LevelLight, LevelBalanced, LevelStrong:
		return true
	}
	return false
}
func (s Status) Terminal() bool       { return s == StatusArchived }
func (s Status) CanPublish() bool     { return s == StatusApproved || s == StatusPublished }
func (r Record) Clone() Record        { return r }
func (r *Record) Touch(now time.Time) { r.UpdatedAt = now }
func (r *Record) Advance(s Status, now time.Time) error {
	if r.Status == StatusArchived {
		return errors.New("record archived")
	}
	r.Status = s
	r.Touch(now)
	return nil
}
func NewRecord(id, name string, level Level, now time.Time) Record {
	return Record{ID: id, Name: name, Level: level, Status: StatusDraft, CreatedAt: now, UpdatedAt: now, Version: 1}
}
func NewAudit(id, record, action, actor, detail string, now time.Time) AuditEvent {
	return AuditEvent{ID: id, RecordID: record, Action: action, Actor: actor, Detail: detail, At: now}
}
func (w Workflow) Valid() error {
	if w.ID == "" {
		return errors.New("workflow id required")
	}
	if len(w.RecordIDs) == 0 {
		return errors.New("records required")
	}
	if w.Stage == "" {
		return errors.New("stage required")
	}
	return nil
}
func (a Attachment) Valid() error {
	if a.ID == "" || a.RecordID == "" || a.Name == "" {
		return errors.New("attachment metadata required")
	}
	if a.Size < 0 {
		return errors.New("negative size")
	}
	return nil
}
