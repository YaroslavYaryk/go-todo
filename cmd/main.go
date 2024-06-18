package main

import (
	"log"
	"simpleRestApi/internal/handler"
	"simpleRestApi/internal/repository"
	"simpleRestApi/internal/service"
	"simpleRestApi/pkg/server"
)

func main() {

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)

	log.Fatal(srv.Run("8080", handlers.InitRoutes()))
}
