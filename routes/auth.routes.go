package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/controllers"
	"github.com/vuongtruongson99/ocr_project/middleware"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.GET("/register", rc.authController.ShowSignUp)
	router.POST("/register", rc.authController.SignUpUser)

	router.GET("/login", rc.authController.ShowSignIn)
	router.POST("/login", rc.authController.SignInUser)

	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)

	router.GET("/text-to-image", middleware.DeserializeUser(), rc.authController.ShowMainTTI)
	router.POST("/text-to-image", middleware.DeserializeUser(), rc.authController.RequestImage)
}
