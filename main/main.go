package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
	"todoproject/api/todo"
	"todoproject/api/users"
	"todoproject/api/util"
	"todoproject/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	fmt.Printf("Config: %v", config)

	// Connect to Postgres
	client, errConnect := db.NewClient(config)
	if errConnect != nil {
		logger.Errorf("failed to connect. due to error: %v", errConnect)
		return
	}

	// Setup Gin
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
	})
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://149.154.65.144"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           time.Duration(viper.GetInt(util.ConfigPath(util.Token, "age"))),
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
	var path string

	if strings.Compare(runtime.GOOS, "linux") == 0 {
		path = "../config"
	} else {
		path = "./config"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("failed read config file. due to error: %v", err)
		return
	}
}
