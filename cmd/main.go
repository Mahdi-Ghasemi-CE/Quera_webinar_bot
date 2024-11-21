package main

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/persistence/database"
	"Quera_webinar_bot/internal/persistence/migrations"
	"Quera_webinar_bot/internal/telegram"
	"Quera_webinar_bot/usecase"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Running the Quera_webinar_bot . . . \n")

	cfg := config.GetConfig()
	fmt.Println("Config file was imported : ", cfg)

	err := database.SetupDb(cfg)
	defer database.CloseDb()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("Postgres was imported ")

	migrations.UpInit()
	fmt.Println("UpInit was migrated ")

	// Setup Telegram bot
	err = telegram.SetupBot(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	defer telegram.CloseBot()

	bot := telegram.GetBot()
	usecase.BotMain(bot)
}
