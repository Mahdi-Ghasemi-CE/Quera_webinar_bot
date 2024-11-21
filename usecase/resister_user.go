package usecase

import (
	"Quera_webinar_bot/internal/telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func BotMain(bot *tgbotapi.BotAPI) {
	// Start receiving updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				err := telegram.SendMessage(update.Message.Chat.ID, "Welcome to Quera Webinar Bot!")
				if err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			}
		}
	}
}
