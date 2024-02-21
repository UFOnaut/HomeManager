package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App              App
		Db               Db
		Jwt              Jwt
		EmailCredentials EmailCredentials
		Endpoint         Endpoint
	}

	App struct {
		Port int
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}

	Jwt struct {
		SecretKey string
	}

	EmailCredentials struct {
		Email    string
		Password string
	}

	Endpoint struct {
		BaseUrl string
	}
)

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}

	return Config{
		App: App{
			Port: viper.GetInt("app.server.port"),
		},
		Db: Db{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			SSLMode:  viper.GetString("database.sslmode"),
			TimeZone: viper.GetString("database.timezone"),
		},
		Jwt: Jwt{
			SecretKey: viper.GetString("jwt_secret_key"),
		},
		EmailCredentials: EmailCredentials{
			Email:    viper.GetString("email_credentials.email"),
			Password: viper.GetString("email_credentials.password"),
		},
		Endpoint: Endpoint{
			BaseUrl: viper.GetString("endpoint.base_url"),
		},
	}
}
