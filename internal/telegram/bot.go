package telegram

import (
	"Quera_webinar_bot/config"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var botClient *tgbotapi.BotAPI

// SetupBot initializes the Telegram bot client
func SetupBot(cfg *config.Config) error {
	var err error
	botClient, err = tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", botClient.Self.UserName)
	return nil
}

// GetBot returns the initialized bot client
func GetBot() *tgbotapi.BotAPI {
	return botClient
}

// CloseBot performs cleanup (optional, for standardization)
func CloseBot() {
	// Add any cleanup logic if necessary
	log.Println("Bot client cleanup completed.")
}

// SendMessage sends a message to a specific chat
func SendMessage(chatID int64, text string) error {
	if botClient == nil {
		return ErrBotNotInitialized
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := botClient.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	return nil
}

// ErrBotNotInitialized is an error returned when the bot is not initialized
var ErrBotNotInitialized = fmt.Errorf("Telegram bot client not initialized")
