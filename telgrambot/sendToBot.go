package telgrambot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}
}

func SendFileToTelegram(botToken, chatID, filePath string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	chatIDInt64, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocument(chatIDInt64, tgbotapi.FilePath(filePath))

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	log.Printf("File sent to Telegram: %s", filePath)
	return nil
}
