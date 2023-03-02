package main

import (
	"fmt"
	"log"

	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/config"
)

func init() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	config.ConnectToDB(&configs)
}

func main() {
	config.DB.AutoMigrate(&models.Comment{}, &models.Movie{}, &models.Character{})
	fmt.Println("Migration completed")
}
