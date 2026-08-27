package bambu

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSeenMissingFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	s, err := LoadSeen(dir)
	if err != nil {
		t.Fatal(err)
	}
	if s.Has("video.avi") {
		t.Fatal("empty seen should not have entries")
	}
}

func TestSeenAddPersists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	s, err := LoadSeen(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add("video.avi"); err != nil {
		t.Fatal(err)
	}
	if !s.Has("video.avi") {
		t.Fatal("expected Has after Add")
	}

	data, err := os.ReadFile(filepath.Join(dir, seenFileName))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "video.avi\n" {
		t.Fatalf("file = %q", data)
	}

	s2, err := LoadSeen(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !s2.Has("video.avi") {
		t.Fatal("reloaded seen missing entry")
	}
}
