package main

import (
	"fmt"
	"log"

	"github.com/vuongtruongson99/ocr_project/initializers"
	"github.com/vuongtruongson99/ocr_project/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{})
	fmt.Println("? Migration complete")
}
