package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vuongtruongson99/ocr_project/initializers"
	"github.com/vuongtruongson99/ocr_project/models"
	"github.com/vuongtruongson99/ocr_project/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// Show form SignUp
func (ac *AuthController) ShowSignUp(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{},
	)
}

// Show SignIn form
func (ac *AuthController) ShowSignIn(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signin.html",
		gin.H{},
	)
}

// Show Text-to-Image form
func (ac *AuthController) ShowMainTTI(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"tti.html",
		gin.H{},
	)
}

// SignUp User
func (ac *AuthController) SignUpUser(c *gin.Context) {
	var payload *models.SignUpInput

	if err := c.Bind(&payload); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if payload.Password != payload.PasswordConfirm {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"status":  "fail",
			"message": "Passwords do not match",
		})
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.HTML(http.StatusBadGateway, "signup.html", gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}

	now := time.Now()
	newUsers := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "user",
		Verified:  true,
		Provider:  "local",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUsers)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		c.HTML(http.StatusConflict, "signup.html", gin.H{
			"status":  "fail",
			"message": "User with that email already exists",
		})
		return
	} else if result.Error != nil {
		c.HTML(http.StatusBadGateway, "signup.html", gin.H{
			"status":  "error",
			"message": "Something bad happened",
		})
		return
	}

	c.HTML(http.StatusCreated, "signup.html", gin.H{
		"status":  "success",
		"message": "Successful created!"})
}

// Login User
func (ac *AuthController) SignInUser(c *gin.Context) {
	var payload *models.SignInInput

	if err := c.Bind(&payload); err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))

	// Check email
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"status":  "fail",
			"message": "Invalid email or password",
		})
		return
	}

	// Check password
	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"status":  "fail",
			"message": "Invalid email or password",
		})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Token
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	c.Redirect(http.StatusFound, "/api/auth/text-to-image")
}

// Refresh access token
func (ac *AuthController) RefreshAccessToken(c *gin.Context) {
	message := "could not refresh access token"

	cookie, err := c.Cookie("refresh_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": message,
		})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_id", "", -1, "/", "localhost", false, false)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"status": "success",
	})
}

// Send request to HF Inference API
func (ac *AuthController) RequestImage(c *gin.Context) {
	var payload *models.GenerateImage

	if err := c.Bind(&payload); err != nil {
		c.HTML(http.StatusBadRequest, "tti.html", gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	config, _ := initializers.LoadConfig(".")
	APIURL := "https://api-inference.huggingface.co/models/" + payload.Model

	var images []string
	imageBytes, err := utils.Query(map[string]interface{}{
		"inputs": payload.Prompt,
	}, APIURL, config.HFAPIToken)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	img2 := base64.StdEncoding.EncodeToString(imageBytes)
	images = append(images, img2)

	// os.WriteFile("output/out.png", imageBytes, 0666)

	c.HTML(http.StatusOK, "tti.html", gin.H{
		"images": images,
	})

}
