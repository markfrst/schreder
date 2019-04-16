package main

import (
	"log"
	"net/http"
	"os"
	"io"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	StorePath, Token := LoadEnv()

	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Document == nil { // ignore any non-Message Updates
			continue
		}

		fileConfig := tgbotapi.FileConfig{FileID: update.Message.Document.FileID}

		file, _ := bot.GetFile(fileConfig)
		fileName := update.Message.Document.FileName
		link := file.Link(Token)

		log.Printf("FileID: %s", file.FileID)
		log.Printf("Link to file: %s", link)

		if err := DownloadFile(StorePath + fileName, link); err != nil {
			panic(err)
		}
	}
}

// DownloadFile download file by link
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
			return err
	}

	log.Printf("resp: %+v", resp)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
			return err
	}
	log.Printf("out: %+v", out)

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// LoadEnv load env vars from .env
func LoadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("STORE_PATH"), os.Getenv("TELEGRAM_TOKEN")
}
