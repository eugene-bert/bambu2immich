package bambu

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPConfig struct {
	IP          string
	AccessCode  string
	DownloadDir string
	FilePrefix  string
}

func DownloadNewTimelapses(cfg FTPConfig, seen map[string]bool) ([]string, error) {
	addr := fmt.Sprintf("%s:990", cfg.IP)
	conn, err := ftp.Dial(addr,
		ftp.DialWithTLS(&tls.Config{InsecureSkipVerify: true}),
		ftp.DialWithTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("FTP dial: %w", err)
	}
	defer conn.Quit()

	if err := conn.Login("bblp", cfg.AccessCode); err != nil {
		return nil, fmt.Errorf("FTP login: %w", err)
	}

	entries, err := conn.List("/timelapse")
	if err != nil {
		return nil, fmt.Errorf("FTP list: %w", err)
	}

	var downloaded []string
	for _, e := range entries {
		if e.Type != ftp.EntryTypeFile || filepath.Ext(e.Name) != ".avi" || seen[e.Name] {
			continue
		}

		log.Printf("downloading timelapse: %s (%d bytes)", e.Name, e.Size)

		resp, err := conn.Retr("/timelapse/" + e.Name)
		if err != nil {
			log.Printf("FTP retr %s: %v", e.Name, err)
			continue
		}

		if err := os.MkdirAll(cfg.DownloadDir, 0o755); err != nil {
			resp.Close()
			return downloaded, fmt.Errorf("mkdir: %w", err)
		}

		localName := e.Name
		if cfg.FilePrefix != "" {
			localName = cfg.FilePrefix + strings.TrimPrefix(e.Name, "video")
		}
		outPath := filepath.Join(cfg.DownloadDir, localName)
		f, err := os.Create(outPath)
		if err != nil {
			resp.Close()
			return downloaded, fmt.Errorf("create file: %w", err)
		}

		_, copyErr := io.Copy(f, resp)
		f.Close()
		resp.Close()

		if copyErr != nil {
			os.Remove(outPath)
			log.Printf("download %s failed: %v", e.Name, copyErr)
			continue
		}

		seen[e.Name] = true
		downloaded = append(downloaded, outPath)
		log.Printf("downloaded: %s", outPath)
	}

	return downloaded, nil
}
