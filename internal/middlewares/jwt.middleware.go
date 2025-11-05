package middlewares

import (
	"mmgrapp/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// memeriksa token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		// Validasi Header
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		// Validasi Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			ctx.Abort()
			return
		}

		// Validasi JWT
		tokenString := parts[1]
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired token"})
			ctx.Abort()
			return
		}

		// set UserID to context
		ctx.Set("user_id", claims.UserID)

		// lanjut ke next handler
		ctx.Next()
	}
}
