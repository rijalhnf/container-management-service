package middleware

import (
	"net/http"
	"strings"

	"go-gin-postgre-crud/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT Bearer token in the Authorization header.
// On success, it sets "userID", "email", and "role" into the gin.Context.
// On failure, it aborts the request with HTTP 401.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be: Bearer <token>",
			})
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Inject user info into context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
