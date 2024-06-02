package main

import (
	"context"
	"time"
	"todoapi/internal/config"
	"todoapi/internal/handler"
	"todoapi/internal/infrastructure"
	"todoapi/internal/repository"

	_ "todoapi/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title To do api description
// @version 2.0
// @description http service
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("viper OK")

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	pgConfig, err := config.GetDBConfig()

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("got bd config")

	db, err := infrastructure.SetUpPostgresDatabase(ctx, pgConfig)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("connected to db")

	taskRepository := repository.NewTaskRepository(db)
	taskHandler := handler.NewTaskHandler(taskRepository)

	taskHandler.InitRoutes(app)

	port := viper.GetString("http.port")
	if err = app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

	log.Info("server is running")
}

func SetupViper() error {

	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
