package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eugene-bert/bambu2immich/bambu"
	"github.com/eugene-bert/bambu2immich/immich"
	"github.com/eugene-bert/bambu2immich/notify"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	mqttCfg := bambu.MQTTConfig{
		IP:         os.Getenv("BAMBU_IP"),
		AccessCode: os.Getenv("BAMBU_ACCESS_CODE"),
		Serial:     os.Getenv("BAMBU_SERIAL"),
	}
	printerName := envOr("BAMBU_NAME", "bambu-a1")
	ftpCfg := bambu.FTPConfig{
		IP:          mqttCfg.IP,
		AccessCode:  mqttCfg.AccessCode,
		DownloadDir: envOr("DOWNLOAD_DIR", "/data/timelapses"),
		FilePrefix:  printerName,
	}
	immichCfg := immich.Config{
		URL:    os.Getenv("IMMICH_URL"),
		APIKey: os.Getenv("IMMICH_API_KEY"),
	}
	telegramCfg := notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
	keepLocal := strings.EqualFold(envOr("KEEP_LOCAL", "false"), "true")

	if mqttCfg.IP == "" || mqttCfg.AccessCode == "" || mqttCfg.Serial == "" {
		log.Fatal("BAMBU_IP, BAMBU_ACCESS_CODE, and BAMBU_SERIAL are required")
	}
	if immichCfg.URL == "" || immichCfg.APIKey == "" {
		log.Fatal("IMMICH_URL and IMMICH_API_KEY are required")
	}

	seen := make(map[string]bool)

	onFinish := func() {
		log.Println("waiting 15s for printer to finalize timelapse...")
		time.Sleep(15 * time.Second)

		paths, err := bambu.DownloadNewTimelapses(ftpCfg, seen)
		if err != nil {
			log.Printf("download error: %v", err)
		}
		if len(paths) == 0 {
			log.Println("no new timelapses found")
			return
		}

		for _, path := range paths {
			desc := fmt.Sprintf("3D print timelapse from %s", printerName)
			if err := immich.Upload(immichCfg, path, desc); err != nil {
				log.Printf("upload error: %v", err)
				continue
			}

			if telegramCfg.Enabled() {
				telegramCfg.Send(path)
			}

			if !keepLocal {
				os.Remove(path)
				log.Printf("cleaned up: %s", path)
			}
		}
	}

	if err := bambu.Listen(mqttCfg, onFinish); err != nil {
		log.Fatal(err)
	}

	log.Println("bambu2immich running — waiting for prints to finish...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}
