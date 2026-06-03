// Package store handles persistence of snapshots to disk.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/DucChau/driftmap/internal/snapshot"
)

const dirName = ".driftmap"

// Store manages snapshot files on disk.
type Store struct {
	dir string
}

// Open returns a Store backed by ~/.driftmap, creating the directory if needed.
func Open() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("finding home dir: %w", err)
	}
	dir := filepath.Join(home, dirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("creating store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot to disk as a JSON file.
func (s *Store) Save(snap *snapshot.Snapshot) error {
	filename := filepath.Join(s.dir, snap.ID+".json")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil // idempotent
		}
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// List returns all snapshots sorted by capture time (oldest first).
func (s *Store) List() ([]*snapshot.Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var snaps []*snapshot.Snapshot
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		snap, err := s.load(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		snaps = append(snaps, snap)
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].CapturedAt.Before(snaps[j].CapturedAt)
	})

	return snaps, nil
}

// Resolve finds a snapshot by full ID or short prefix or label.
func (s *Store) Resolve(ref string) (*snapshot.Snapshot, error) {
	snaps, err := s.List()
	if err != nil {
		return nil, err
	}

	// Try exact ID or prefix match first.
	for _, snap := range snaps {
		if snap.ID == ref || strings.HasPrefix(snap.ID, ref) {
			return snap, nil
		}
	}

	// Try label match.
	for _, snap := range snaps {
		if snap.Label == ref {
			return snap, nil
		}
	}

	// Try special aliases: "latest", "oldest".
	if ref == "latest" && len(snaps) > 0 {
		return snaps[len(snaps)-1], nil
	}
	if ref == "oldest" && len(snaps) > 0 {
		return snaps[0], nil
	}

	// Try time-offset: "@-1" means 1 snapshot before latest.
	if strings.HasPrefix(ref, "@-") {
		var offset int
		if _, err := fmt.Sscanf(ref[1:], "%d", &offset); err == nil && offset < 0 {
			idx := len(snaps) + offset
			if idx >= 0 && idx < len(snaps) {
				return snaps[idx], nil
			}
		}
	}

	return nil, fmt.Errorf("snapshot %q not found", ref)
}

func (s *Store) load(path string) (*snapshot.Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var snap snapshot.Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	if snap.CapturedAt.IsZero() {
		snap.CapturedAt = time.Now()
	}
	return &snap, nil
}
