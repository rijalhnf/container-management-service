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
// POST /api/auth/register
// Body: { "email": "...", "password": "..." }
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
// POST /api/auth/login
// Body: { "email": "...", "password": "..." }
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
// GET /api/auth/me  (protected route)
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
