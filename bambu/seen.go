package bambu

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const seenFileName = ".seen"

// Seen is a persisted set of remote timelapse filenames already downloaded.
type Seen struct {
	mu    sync.Mutex
	names map[string]struct{}
	path  string
}

func LoadSeen(dir string) (*Seen, error) {
	s := &Seen{
		names: make(map[string]struct{}),
		path:  filepath.Join(dir, seenFileName),
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s.names[line] = struct{}{}
	}
	return s, nil
}

func (s *Seen) Has(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.names[name]
	return ok
}

func (s *Seen) Add(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.names[name]; ok {
		return nil
	}
	s.names[name] = struct{}{}
	return s.flushLocked()
}

func (s *Seen) flushLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	var b strings.Builder
	for name := range s.names {
		b.WriteString(name)
		b.WriteByte('\n')
	}
	return os.WriteFile(s.path, []byte(b.String()), 0o644)
}
