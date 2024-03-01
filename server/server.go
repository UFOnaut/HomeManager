package server

import (
	"fmt"

	"home_manager/config"
	"home_manager/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"home_manager/repositories"
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
	loginHttpHandler := handlers.NewLoginHttpHandler(userRepository)
	registerHttpHandler := handlers.NewRegisterHttpHandler(userRepository)

	// Routers
	s.app.POST("/login", loginHttpHandler.Login)
	s.app.POST("/refresh_token", loginHttpHandler.RefreshToken)
	s.app.POST("/register", registerHttpHandler.Register)
	s.app.GET("/verify", registerHttpHandler.VerifyEmail)
}
