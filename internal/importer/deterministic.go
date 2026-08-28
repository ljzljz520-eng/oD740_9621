package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"example.com/noise-registry/internal/model"
	"fmt"
	"time"
)

func Checksum(row Row) string {
	h := sha256.Sum256([]byte(row.ID + "|" + row.Name + "|" + row.Level))
	return hex.EncodeToString(h[:])
}
func AttachmentFor(row Row, now time.Time) model.Attachment {
	return model.Attachment{ID: row.ID + "-source", RecordID: row.ID, Name: row.Name + ".wav", Checksum: Checksum(row), Size: int64(len(row.Name)), CreatedAt: now}
}
func Describe(r Result) string { return fmt.Sprintf("accepted=%d rejected=%d", r.Accepted, r.Rejected) }
func AcceptedIDs(r Result) []string {
	out := []string{}
	for _, x := range r.Records {
		out = append(out, x.ID)
	}
	return out
}
func Rejected(r Result) bool { return r.Rejected > 0 }
