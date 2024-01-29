package main

import (
	"home_manager/config"
	"home_manager/database"
	"home_manager/entities"
)

func main() {
	cfg := config.GetConfig()
	db := database.NewPostgresDatabase(&cfg)
	migrate(db)

}

func migrate(db database.Database) {
	db.GetDb().Migrator().CreateTable(&entities.User{})
	db.GetDb().Create(&entities.User{
		Id:       3,
		Email:    "test@gmail.com",
		Name:     "Illia",
		GroupIds: entities.GroupIds{1, 2, 3},
	})
}
