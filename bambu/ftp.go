package bambu

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPConfig struct {
	IP          string
	AccessCode  string
	DownloadDir string
	FilePrefix  string
}

func DownloadNewTimelapses(cfg FTPConfig, seen *Seen) ([]string, error) {
	addr := fmt.Sprintf("%s:990", cfg.IP)
	conn, err := ftp.Dial(addr,
		ftp.DialWithTLS(printerTLS()),
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
		if e.Type != ftp.EntryTypeFile {
			continue
		}
		name, err := sanitizeRemoteName(e.Name)
		if err != nil {
			continue
		}
		if seen != nil && seen.Has(name) {
			continue
		}

		log.Printf("downloading timelapse: %s (%d bytes)", name, e.Size)

		resp, err := conn.Retr("/timelapse/" + name)
		if err != nil {
			log.Printf("FTP retr %s: %v", name, err)
			continue
		}

		if err := os.MkdirAll(cfg.DownloadDir, 0o755); err != nil {
			resp.Close()
			return downloaded, fmt.Errorf("mkdir: %w", err)
		}

		outPath, err := confinedPath(cfg.DownloadDir, localFilename(name, cfg.FilePrefix))
		if err != nil {
			resp.Close()
			log.Printf("skip %s: %v", name, err)
			continue
		}

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
			log.Printf("download %s failed: %v", name, copyErr)
			continue
		}

		if !e.Time.IsZero() {
			if err := os.Chtimes(outPath, e.Time, e.Time); err != nil {
				log.Printf("chtimes %s: %v", outPath, err)
			}
		}

		if seen != nil {
			if err := seen.Add(name); err != nil {
				log.Printf("persist seen %s: %v", name, err)
			}
		}
		downloaded = append(downloaded, outPath)
		log.Printf("downloaded: %s", outPath)
	}

	return downloaded, nil
}
