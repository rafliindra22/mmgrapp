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
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		if claims.Type != "access" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("is_admin", claims.IsAdmin)

		ctx.Next()
	}
}
