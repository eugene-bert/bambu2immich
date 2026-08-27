package immich

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadRejectsRedirect(t *testing.T) {
	t.Parallel()
	var sawAPIKey bool
	redir := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "" {
			sawAPIKey = true
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"stolen"}`))
	}))
	t.Cleanup(redir.Close)

	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redir.URL, http.StatusFound)
	}))
	t.Cleanup(origin.Close)

	path := filepath.Join(t.TempDir(), "clip.avi")
	if err := os.WriteFile(path, []byte("avi"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Upload(Config{URL: origin.URL, APIKey: "secret-key"}, path, "desc")
	if err == nil {
		t.Fatal("expected redirect to fail")
	}
	if !strings.Contains(err.Error(), "302") {
		t.Fatalf("error = %v", err)
	}
	if sawAPIKey {
		t.Fatal("API key was sent to redirect target")
	}
}

func TestUploadOK(t *testing.T) {
	t.Parallel()
	var apiKey, ct string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey = r.Header.Get("x-api-key")
		ct = r.Header.Get("Content-Type")
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"abc","duplicate":false}`))
	}))
	t.Cleanup(srv.Close)

	path := filepath.Join(t.TempDir(), "clip.avi")
	if err := os.WriteFile(path, []byte("avi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Upload(Config{URL: srv.URL + "/", APIKey: "k"}, path, "timelapse"); err != nil {
		t.Fatal(err)
	}
	if apiKey != "k" {
		t.Errorf("api key = %q", apiKey)
	}
	if !strings.HasPrefix(ct, "multipart/form-data") {
		t.Errorf("content-type = %q", ct)
	}
}
