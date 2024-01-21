package main

import (
	"home_manager/config"
	"home_manager/database"
	"home_manager/server"
)

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(&cfg)

	server.NewEchoServer(&cfg, db.GetDb()).Start()
}
