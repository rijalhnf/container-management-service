package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"
	"go-gin-postgre-crud/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account.
// @Summary      Register a new user
// @Description  Create a new user account with email and password. Password must be at least 6 characters.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "User registration payload"
// @Success      201 {object} map[string]interface{} "User registered successfully"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Email is already registered"
// @Router       /api/auth/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existing models.User
	if result := config.DB.Where("email = ?", req.Email).First(&existing); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	// Hash the password with bcrypt (cost factor 12)
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     "user",
	}

	if result := config.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login authenticates a user and returns a signed JWT token.
// @Summary      Login user
// @Description  Authenticate with email and password to receive a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login credentials"
// @Success      200 {object} map[string]interface{} "Login success, returns token"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      401 {object} map[string]string "Invalid email or password"
// @Router       /api/auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up the user by email
	var user models.User
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		// Use a generic message to avoid leaking whether the email exists
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Me returns the currently authenticated user's info extracted from JWT.
// @Summary      Get current user
// @Description  Returns the currently authenticated user's info extracted from JWT claims.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Current user info"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /api/auth/me [get]
func Me(c *gin.Context) {
	userID, _ := c.Get("userID")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}
