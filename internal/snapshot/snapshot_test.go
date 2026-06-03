package snapshot_test

import (
	"testing"

	"github.com/DucChau/driftmap/internal/snapshot"
)

func TestNew_ParsesKeyValue(t *testing.T) {
	snap := snapshot.New("test", []string{"FOO=bar", "BAZ=qux"})

	if snap.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", snap.Vars["FOO"])
	}
	if snap.Vars["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", snap.Vars["BAZ"])
	}
	if snap.Label != "test" {
		t.Errorf("expected label=test, got %q", snap.Label)
	}
}

func TestNew_HandlesValueWithEquals(t *testing.T) {
	snap := snapshot.New("", []string{"URL=postgres://user:pass@host/db?sslmode=disable"})
	expected := "postgres://user:pass@host/db?sslmode=disable"
	if snap.Vars["URL"] != expected {
		t.Errorf("expected %q, got %q", expected, snap.Vars["URL"])
	}
}

func TestShortID(t *testing.T) {
	snap := snapshot.New("", []string{"X=1"})
	short := snap.ShortID()
	if len(short) != 8 {
		t.Errorf("expected 8-char short ID, got %d: %q", len(short), short)
	}
	if snap.ID[:8] != short {
		t.Error("ShortID should be first 8 chars of ID")
	}
}
