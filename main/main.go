package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"todoproject/api/todo"
	"todoproject/api/users"
	"todoproject/api/util"
	"todoproject/db"
)

func main() {
	runApp()
}

func runApp() {
	// Configuration Viper
	ConfigureViper()

	logger := logrus.StandardLogger()

	// Setup Database Config
	config := &db.Config{
		Username:     viper.GetString(util.ConfigPath(util.Postgres, "username")),
		Password:     viper.GetString(util.ConfigPath(util.Postgres, "password")),
		Host:         viper.GetString(util.ConfigPath(util.Postgres, "host")),
		Port:         viper.GetString(util.ConfigPath(util.Postgres, "port")),
		Database:     viper.GetString(util.ConfigPath(util.Postgres, "database")),
		TokenExpires: viper.GetDuration(util.ConfigPath(util.Token, "expires_in")),
		TokenSecret:  viper.GetString(util.ConfigPath(util.Token, "secret_key")),
		TokenMaxAge:  viper.GetInt(util.ConfigPath(util.Token, "age")),
	}

	// Connect to Postgres
	client, errConnect := db.NewClient(config)
	if errConnect != nil {
		logger.Errorf("failed to connect. due to error: %v", errConnect)
		return
	}

	// Setup Gin
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:           []string{"http://localhost"},
		AllowOriginFunc:        nil,
		AllowMethods:           nil,
		AllowHeaders:           nil,
		AllowCredentials:       false,
		ExposeHeaders:          nil,
		MaxAge:                 0,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}))

	// init users storageUsers
	storageUsers := users.NewStorage(client, logger)
	// init users controller
	userHandler := users.NewHandler(storageUsers, config, logger)
	userHandler.InitUserHandler(server)

	// init storage todos
	storageTodos := todo.NewStorage(client, logger)
	// init storage controller
	todosHandler := todo.NewHandler(storageTodos, userHandler, logger)
	todosHandler.InitTodoHandler(server)

	log.Fatalln(server.Run(viper.GetString(util.ConfigPath(util.Server, "port"))))
}

func ConfigureViper() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("failed read config file. due to error: %v", err)
		return
	}
}
