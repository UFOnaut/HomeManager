package server

import (
	"fmt"

	"home_manager/config"

	"home_manager/handlers"
	"home_manager/repositories"
	"home_manager/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server interface {
	Start()
}

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
	s.initializeUserHttpHandler()
	s.app.Use(middleware.Logger())

	serverUrl := fmt.Sprintf(":%d", s.cfg.App.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *echoServer) initializeUserHttpHandler() {
	// Initialize all layers
	userRepository := repositories.NewUserRepository(s.db)
	userUsecase := usecases.NewUserUsecase(userRepository)
	userHttpHandler := handlers.NewUserHttpHandler(userUsecase)

	// Routers
	s.app.POST("/login", userHttpHandler.Login)
}
