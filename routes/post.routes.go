package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/controllers"
	"github.com/vuongtruongson99/ocr_project/middleware"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewRoutePostController(postController controllers.PostController) PostRouteController {
	return PostRouteController{postController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {
	router := rg.Group("posts")
	router.Use(middleware.DeserializeUser())
	router.POST("/", pc.postController.CreatePost) // Create new post
	router.GET("/", pc.postController.FindPosts)   // Get all posts

	router.GET("/:postId", pc.postController.FindPostById)
	router.PUT("/:postId", pc.postController.UpdatePost)
	router.DELETE("/:postId", pc.postController.DeletePost)

}
