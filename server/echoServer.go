package server

import (
	"fmt"

	"home_manager/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type echoServer struct {
	app *echo.Echo
	db  *gorm.DB
	cfg *config.Config
}

func NewEchoServer(cfg *config.Config, db *gorm.DB) Server {
	return &echoServer{
		app: echo.New(),
		db:  db,
		cfg: cfg,
	}
}

func (s *echoServer) Start() {
	// Initialzie routers here
	// ...

	s.app.Use(middleware.Logger())

	serverUrl := fmt.Sprintf(":%d", s.cfg.App.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}
