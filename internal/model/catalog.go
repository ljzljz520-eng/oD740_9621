package model

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Catalog struct {
	Records     map[string]Record
	Audits      map[string][]AuditEvent
	Workflows   map[string]Workflow
	Attachments map[string][]Attachment
}

func NewCatalog() Catalog {
	return Catalog{Records: map[string]Record{}, Audits: map[string][]AuditEvent{}, Workflows: map[string]Workflow{}, Attachments: map[string][]Attachment{}}
}
func (c *Catalog) AddRecord(r Record) error {
	if e := r.Validate(); e != nil {
		return e
	}
	c.Records[r.ID] = r
	return nil
}
func (c *Catalog) AddAudit(a AuditEvent) { c.Audits[a.RecordID] = append(c.Audits[a.RecordID], a) }
func (c *Catalog) AddWorkflow(w Workflow) error {
	if e := w.Valid(); e != nil {
		return e
	}
	c.Workflows[w.ID] = w
	return nil
}
func (c *Catalog) AddAttachment(a Attachment) error {
	if e := a.Valid(); e != nil {
		return e
	}
	c.Attachments[a.RecordID] = append(c.Attachments[a.RecordID], a)
	return nil
}
func (c Catalog) Record(id string) (Record, bool)     { r, ok := c.Records[id]; return r, ok }
func (c Catalog) Audit(id string) []AuditEvent        { return append([]AuditEvent{}, c.Audits[id]...) }
func (c Catalog) Workflow(id string) (Workflow, bool) { w, ok := c.Workflows[id]; return w, ok }
func (c Catalog) Attachment(id string) []Attachment {
	return append([]Attachment{}, c.Attachments[id]...)
}
func (c Catalog) Names() []string {
	out := []string{}
	for _, r := range c.Records {
		out = append(out, r.Name)
	}
	sort.Strings(out)
	return out
}
func (c Catalog) ByStatus(s Status) []Record {
	out := []Record{}
	for _, r := range c.Records {
		if r.Status == s {
			out = append(out, r)
		}
	}
	return SortRecords(out)
}
func (c Catalog) Search(q string) []Record {
	out := []Record{}
	for _, r := range c.Records {
		if strings.Contains(strings.ToLower(r.Name), strings.ToLower(q)) {
			out = append(out, r)
		}
	}
	return SortRecords(out)
}
func (c Catalog) Summary() string {
	return fmt.Sprintf("records=%d workflows=%d", len(c.Records), len(c.Workflows))
}
func (c *Catalog) Touch(id string, now time.Time) error {
	r, ok := c.Records[id]
	if !ok {
		return fmt.Errorf("missing %s", id)
	}
	r.Touch(now)
	c.Records[id] = r
	return nil
}
func (c Catalog) PublishedNames() []string { rs := c.ByStatus(StatusPublished); return Names(rs) }
func Names(rs []Record) []string {
	out := []string{}
	for _, r := range rs {
		out = append(out, r.Name)
	}
	return out
}
func (c Catalog) Valid() error {
	for _, r := range c.Records {
		if e := r.Validate(); e != nil {
			return e
		}
	}
	return nil
}
