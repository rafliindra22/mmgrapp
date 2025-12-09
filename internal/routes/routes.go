package routes

import (
	config "mmgrapp/internal/configs"
	"mmgrapp/internal/handlers"
	"mmgrapp/internal/middlewares"
	"mmgrapp/internal/repositories"
	"mmgrapp/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mendaftarkan semua route ke server
func SetupRoutes(r *gin.Engine) {
	db := config.DB
	// ================= USER MODULE =================
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOTPRepository(db)
	userService := services.NewUserService(userRepo, otpRepo)
	userHandler := handlers.NewUserHandler(userService)

	// ================= AUTH MODULE =================
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, userRepo, otpRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Test endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		// user module
		auth.POST("/register", userHandler.Register)
		auth.POST("/verify-email", userHandler.VerifyEmail)
		auth.POST("/resend-otp", userHandler.ResendOTP)

		// auth module
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-pass", authHandler.ForgotPassword)
		auth.POST("/reset-pass", authHandler.ResetPassword)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.RevokeToken)

		profile := api.Group("/profile")
		// profile module
		profile.GET("/my-detail/:id", middlewares.JWTAuthMiddleware(), userHandler.MyDetail)
	}
}
