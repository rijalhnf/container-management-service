package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // "-" hides password from JSON output
	Role      string    `gorm:"type:varchar(50);default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest is the expected request body for registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest is the expected request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after a successful login
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
	User      struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}
