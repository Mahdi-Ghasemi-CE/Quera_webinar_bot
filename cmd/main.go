package main

import (
	"Quera_webinar_bot/config"
	"fmt"
)

func main() {
	fmt.Println("Start")

	cfg := config.GetConfig()
	fmt.Println("Config file was imported : ", cfg)
}
