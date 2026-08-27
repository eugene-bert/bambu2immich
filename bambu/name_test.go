package bambu

import (
	"path/filepath"
	"testing"
)

func TestSanitizeRemoteName(t *testing.T) {
	t.Parallel()
	ok, err := sanitizeRemoteName("video.avi")
	if err != nil || ok != "video.avi" {
		t.Fatalf("got %q, %v", ok, err)
	}

	rejects := []string{
		"",
		"video.mp4",
		"../video.avi",
		"..",
		"foo/../video.avi",
		"/timelapse/video.avi",
		`..\video.avi`,
		"foo..avi",
		".avi",
	}
	for _, name := range rejects {
		if _, err := sanitizeRemoteName(name); err == nil {
			t.Errorf("accepted %q", name)
		}
	}
}

func TestConfinedPath(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "tl")
	out, err := confinedPath(dir, "bambu-a1.avi")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(out) != filepath.Clean(dir) {
		t.Fatalf("dir = %q", filepath.Dir(out))
	}

	if _, err := confinedPath(dir, "../escape.avi"); err == nil {
		t.Fatal("expected escape to fail")
	}
	if _, err := confinedPath(dir, "a/b.avi"); err == nil {
		t.Fatal("expected nested name to fail")
	}
}

func TestLocalFilename(t *testing.T) {
	t.Parallel()
	got := localFilename("video123.avi", "bambu-a1")
	if got != "bambu-a1123.avi" {
		t.Fatalf("got %q", got)
	}
	got = localFilename("clip.avi", "")
	if got != "clip.avi" {
		t.Fatalf("got %q", got)
	}
}
