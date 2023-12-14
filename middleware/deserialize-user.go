package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/initializers"
	"github.com/vuongtruongson99/ocr_project/models"
	"github.com/vuongtruongson99/ocr_project/utils"
)

func DeserializeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var access_token string
		cookie, err := c.Cookie("access_token")

		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			c.HTML(http.StatusUnauthorized, "home.html", gin.H{"status": "fail", "message": "You are not logged in"})
			c.Abort()
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)

		if err != nil {
			c.HTML(http.StatusUnauthorized, "home.html", gin.H{"status": "fail", "message": err.Error()})
			c.Abort()
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			c.HTML(http.StatusForbidden, "home.html", gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
