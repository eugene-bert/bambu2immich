package immich

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

var httpClient = &http.Client{
	Timeout: 15 * time.Minute,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
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

	createdAt := stat.ModTime().UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	ts := createdAt.Format("2006-01-02T15:04:05.000Z")

	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	go func() {
		var pipeErr error
		defer func() {
			w.Close()
			pw.CloseWithError(pipeErr)
		}()

		part, err := w.CreateFormFile("assetData", filepath.Base(filePath))
		if err != nil {
			pipeErr = err
			return
		}
		if _, err := io.Copy(part, f); err != nil {
			pipeErr = err
			return
		}
		if err := w.WriteField("fileCreatedAt", ts); err != nil {
			pipeErr = err
			return
		}
		if err := w.WriteField("fileModifiedAt", ts); err != nil {
			pipeErr = err
			return
		}
		if description != "" {
			if err := w.WriteField("description", description); err != nil {
				pipeErr = err
				return
			}
		}
	}()

	url := strings.TrimRight(cfg.URL, "/") + "/api/assets"
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("x-api-key", cfg.APIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(body))
	}

	var result uploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("uploaded %s but could not parse Immich response: %v", filepath.Base(filePath), err)
		return nil
	}

	if result.Duplicate {
		log.Printf("skipped duplicate: %s", filepath.Base(filePath))
	} else {
		log.Printf("uploaded to Immich: %s (id: %s)", filepath.Base(filePath), result.ID)
	}

	return nil
}
