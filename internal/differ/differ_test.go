package differ_test

import (
	"testing"

	"github.com/DucChau/driftmap/internal/differ"
	"github.com/DucChau/driftmap/internal/snapshot"
)

func makeSnap(vars map[string]string) *snapshot.Snapshot {
	env := make([]string, 0, len(vars))
	for k, v := range vars {
		env = append(env, k+"="+v)
	}
	return snapshot.New("", env)
}

func TestDiff_Added(t *testing.T) {
	a := makeSnap(map[string]string{"FOO": "1"})
	b := makeSnap(map[string]string{"FOO": "1", "BAR": "2"})
	res := differ.Diff(a, b)

	if len(res.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(res.Added))
	}
	if res.Added["BAR"] != "2" {
		t.Errorf("expected BAR=2, got %q", res.Added["BAR"])
	}
	if len(res.Removed) != 0 || len(res.Changed) != 0 {
		t.Error("expected no removed or changed")
	}
}

func TestDiff_Removed(t *testing.T) {
	a := makeSnap(map[string]string{"FOO": "1", "BAR": "2"})
	b := makeSnap(map[string]string{"FOO": "1"})
	res := differ.Diff(a, b)

	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(res.Removed))
	}
	if res.Removed["BAR"] != "2" {
		t.Errorf("expected BAR=2, got %q", res.Removed["BAR"])
	}
}

func TestDiff_Changed(t *testing.T) {
	a := makeSnap(map[string]string{"FOO": "old"})
	b := makeSnap(map[string]string{"FOO": "new"})
	res := differ.Diff(a, b)

	if len(res.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(res.Changed))
	}
	c := res.Changed["FOO"]
	if c.Before != "old" || c.After != "new" {
		t.Errorf("expected old->new, got %q->%q", c.Before, c.After)
	}
}

func TestDiff_Identical(t *testing.T) {
	a := makeSnap(map[string]string{"FOO": "1", "BAR": "2"})
	b := makeSnap(map[string]string{"FOO": "1", "BAR": "2"})
	res := differ.Diff(a, b)

	if len(res.Added)+len(res.Removed)+len(res.Changed) != 0 {
		t.Error("expected no diff for identical snapshots")
	}
}
