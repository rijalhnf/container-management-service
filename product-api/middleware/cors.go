package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware adds CORS headers to every response.
// Allowed origin is controlled by the CORS_ORIGIN env var (default: "*").
// In production, set CORS_ORIGIN to your specific frontend domain.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := os.Getenv("CORS_ORIGIN")
		if origin == "" {
			origin = "*"
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Max-Age", "86400") // Cache preflight for 24h

		// If credentials are needed (cookies / Auth header with specific origin):
		// c.Header("Access-Control-Allow-Credentials", "true")

		// Handle CORS preflight (OPTIONS) — reply 204 and stop
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
