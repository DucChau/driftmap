// Package differ computes the delta between two snapshots.
package differ

import "github.com/DucChau/driftmap/internal/snapshot"

// Change represents a variable whose value changed between two snapshots.
type Change struct {
	Before string
	After  string
}

// Result holds the complete diff between two snapshots.
type Result struct {
	Added   map[string]string // vars present in B but not A
	Removed map[string]string // vars present in A but not B
	Changed map[string]Change // vars present in both but with different values
}

// Diff computes the diff from snapshot a to snapshot b.
func Diff(a, b *snapshot.Snapshot) Result {
	res := Result{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]Change),
	}

	// Check what's in B vs A.
	for k, bv := range b.Vars {
		if av, ok := a.Vars[k]; !ok {
			res.Added[k] = bv
		} else if av != bv {
			res.Changed[k] = Change{Before: av, After: bv}
		}
	}

	// Check what was in A but missing from B.
	for k, av := range a.Vars {
		if _, ok := b.Vars[k]; !ok {
			res.Removed[k] = av
		}
	}

	return res
}
