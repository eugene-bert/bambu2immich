package notify

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

func (c TelegramConfig) Enabled() bool {
	return c.BotToken != "" && c.ChatID != ""
}

func (c TelegramConfig) Send(filePath string) {
	name := filepath.Base(filePath)
	text := fmt.Sprintf("🎬 Timelapse uploaded to Immich: %s", name)

	resp, err := http.PostForm(
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
	resp.Body.Close()
}
