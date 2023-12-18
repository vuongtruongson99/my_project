package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/controllers"
	"github.com/vuongtruongson99/ocr_project/initializers"
	"github.com/vuongtruongson99/ocr_project/models"
	"github.com/vuongtruongson99/ocr_project/routes"
)

var (
	server         *gin.Engine
	AuthController controllers.AuthController
	UserController controllers.UserController
	PostController controllers.PostController

	AuthRouteController routes.AuthRouteController
	UserRouteController routes.UserRouteController
	PostRouteController routes.PostRouteController
)

func showIndexPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"home.html",
		gin.H{
			"title": "Home Page",
		},
	)
}

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
	AuthController = controllers.NewAuthController(initializers.DB)
	UserController = controllers.NewUserController(initializers.DB)
	PostController = controllers.NewPostController(initializers.DB)

	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserRouteController = routes.NewRouteUserController(UserController)
	PostRouteController = routes.NewRoutePostController(PostController)

	server = gin.Default()
	server.LoadHTMLGlob("templates/template/*")
	server.Static("static/", "./templates/static")

	server.GET("/", showIndexPage)
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}
	initializers.DB.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.GET("/sessions/oauth/google", controllers.GoogleOauth)

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
