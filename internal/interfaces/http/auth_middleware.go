package http

import (
	"credit-layer/internal/application"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const appContextKey = "app"

func AuthMiddleware(appUseCase *application.AppUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		app, err := appUseCase.GetAppByAPIKey(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}

		c.Set(appContextKey, app)
		c.Next()
	}
}
