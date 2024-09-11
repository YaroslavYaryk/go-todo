package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"simpleRestApi/internal/handler"
	"simpleRestApi/internal/repository"
	"simpleRestApi/internal/service"
	"simpleRestApi/pkg/psql"
	"simpleRestApi/pkg/server"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatal("error initialization config", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("error loading .env file", err)
	}

	db, err := psql.NewPostgresDB(psql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)

	newHandlers := srv.AllowCORS(handlers.InitRoutes())

	logrus.Fatal(srv.Run(viper.GetString("port"), newHandlers))
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
