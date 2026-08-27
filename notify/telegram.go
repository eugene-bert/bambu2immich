package notify

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

func (c TelegramConfig) Enabled() bool {
	return c.BotToken != "" && c.ChatID != ""
}

func (c TelegramConfig) Send(filePath string) {
	name := filepath.Base(filePath)
	text := fmt.Sprintf("🎬 Timelapse uploaded to Immich: %s", name)

	resp, err := httpClient.PostForm(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.BotToken),
		url.Values{
			"chat_id": {c.ChatID},
			"text":    {text},
		},
	)
	if err != nil {
		log.Printf("telegram send error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		log.Printf("telegram send error: status %d: %s", resp.StatusCode, body)
	}
}
