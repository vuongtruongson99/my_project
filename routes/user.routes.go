package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/controllers"
	"github.com/vuongtruongson99/ocr_project/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middleware.OauthDeserializeUser(), uc.userController.GetMe)
}
