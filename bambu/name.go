package bambu

import (
	"fmt"
	"path/filepath"
	"strings"
)

func sanitizeRemoteName(name string) (string, error) {
	if name == "" || name != filepath.Base(name) {
		return "", fmt.Errorf("unsafe FTP name %q", name)
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return "", fmt.Errorf("unsafe FTP name %q", name)
	}
	ext := filepath.Ext(name)
	if !strings.EqualFold(ext, ".avi") {
		return "", fmt.Errorf("not an .avi: %q", name)
	}
	if strings.TrimSuffix(name, ext) == "" {
		return "", fmt.Errorf("unsafe FTP name %q", name)
	}
	return name, nil
}

func localFilename(remoteName, prefix string) string {
	name := remoteName
	if prefix != "" {
		name = prefix + strings.TrimPrefix(remoteName, "video")
	}
	return filepath.Base(name)
}

func confinedPath(dir, name string) (string, error) {
	if name == "" || name != filepath.Base(name) {
		return "", fmt.Errorf("unsafe local name %q", name)
	}
	dir = filepath.Clean(dir)
	out := filepath.Join(dir, name)
	rel, err := filepath.Rel(dir, out)
	if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", fmt.Errorf("path escapes download dir: %q", name)
	}
	return out, nil
}
