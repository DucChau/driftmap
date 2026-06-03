// Package snapshot defines the Snapshot type and construction helpers.
package snapshot

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Snapshot is an immutable point-in-time capture of environment variables.
type Snapshot struct {
	ID         string            `json:"id"`
	Label      string            `json:"label"`
	CapturedAt time.Time         `json:"captured_at"`
	Vars       map[string]string `json:"vars"`
}

// New creates a Snapshot from a slice of "KEY=VALUE" strings (as returned by os.Environ).
func New(label string, environ []string) *Snapshot {
	vars := make(map[string]string, len(environ))
	for _, e := range environ {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		} else {
			vars[parts[0]] = ""
		}
	}

	now := time.Now().UTC()
	id := generateID(now, vars)

	return &Snapshot{
		ID:         id,
		Label:      label,
		CapturedAt: now,
		Vars:       vars,
	}
}

// ShortID returns the first 8 characters of the snapshot ID.
func (s *Snapshot) ShortID() string {
	if len(s.ID) >= 8 {
		return s.ID[:8]
	}
	return s.ID
}

func generateID(t time.Time, vars map[string]string) string {
	h := sha256.New()
	_, _ = fmt.Fprintf(h, "%d", t.UnixNano())
	for k, v := range vars {
		_, _ = fmt.Fprintf(h, "%s=%s;", k, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
