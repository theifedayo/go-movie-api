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

// Auto-Migrates models in models/ to DB.
func main() {
	config.DB.AutoMigrate(&models.Comment{}, &models.Movie{}, &models.Character{})
	fmt.Println("Migration completed")
}
