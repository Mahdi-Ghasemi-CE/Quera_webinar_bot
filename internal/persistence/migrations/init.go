package migrations

import (
	"Quera_webinar_bot/internal/models"
	"Quera_webinar_bot/internal/persistence/database"
	"gorm.io/gorm"
	"log"
)

const countStarExp = "count(*)"

func UpInit() {
	database := database.GetDb()

	createTables(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// Basic
	tables = addNewTable(database, models.User{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("tables created")
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func DownInit() {
	// nothing
}
