package handlers

import (
	"mmgrapp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(ctx.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil! Silakan verifikasi email Anda.",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// handler verify OTP
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func (h *UserHandler) VerifyEmail(ctx *gin.Context) {
	var req VerifyOTPRequest

	// validate input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call service
	err := h.userService.VerifyEmailOTP(ctx, req.Email, req.OTP)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email berhasil diverifikasi!",
	})
}

type ResendOTPRequest struct {
	Purpose string `json:"purpose" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

func (h *UserHandler) ResendOTP(ctx *gin.Context) {
	var req ResendOTPRequest

	// validate input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call service
	err := h.userService.ResendOTP(ctx, req.Email, req.Purpose)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP berhasil dikirim ulang!",
	})
}

func (h *UserHandler) MyDetail(ctx *gin.Context) {
	// ambil user_id dari token
	tokenUserID := ctx.GetInt("user_id")
	isAdmin := ctx.GetBool("is_admin")

	// ambil user_id dari param URL
	paramID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	// cek hak akses
	if !isAdmin && paramID != tokenUserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Tidak dapat mengakses data user ini"})
		return
	}

	userData, err := h.userService.GetUserByID(ctx, paramID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get User Data berhasil",
		"data":    userData,
	})
}
