package immich

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	URL    string
	APIKey string
}

type uploadResponse struct {
	ID        string `json:"id"`
	Duplicate bool   `json:"duplicate"`
}

func Upload(cfg Config, filePath string, description string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("assetData", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	w.WriteField("fileCreatedAt", now)
	w.WriteField("fileModifiedAt", stat.ModTime().UTC().Format("2006-01-02T15:04:05.000Z"))
	if description != "" {
		w.WriteField("description", description)
	}
	w.Close()

	req, err := http.NewRequest("POST", cfg.URL+"/api/assets", &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("x-api-key", cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
	}

	var result uploadResponse
	json.Unmarshal(body, &result)

	if result.Duplicate {
		log.Printf("skipped duplicate: %s", filepath.Base(filePath))
	} else {
		log.Printf("uploaded to Immich: %s (id: %s)", filepath.Base(filePath), result.ID)
	}

	return nil
}
