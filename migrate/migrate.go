package main

import (
	"log"

	"github.com/vuongtruongson99/ocr_project/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

// func main() {
// 	initializers.DB.AutoMigrate(&models.User{})
// 	fmt.Println("? Migration complete")
// }
