package main

import (
	"github.com/joho/godotenv"
	"github.com/k4zb3k/todo"
	"github.com/k4zb3k/todo/pkg/handler"
	"github.com/k4zb3k/todo/pkg/repository"
	"github.com/k4zb3k/todo/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initing config %s", err.Error())
	}

	if err := godotenv.Load("/home/k4zb3k/Desktop/todo/.env"); err != nil {
		logrus.Fatalf("error loading env vars %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initing db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error while running http server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("/home/k4zb3k/Desktop/todo/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
