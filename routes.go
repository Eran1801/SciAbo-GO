package main

import (
	"sci-abo-go/middleware"
	"sci-abo-go/requests"

	"github.com/gin-gonic/gin"
)

func InitializerRoutes() *gin.Engine {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// auth routes
	router.POST("/api/auth/register", requests.CreateUser)
	router.POST("/api/auth/login", requests.Login)
	router.POST("/api/auth/change_password", middleware.RequiredAuth, requests.ChangePassword)

	// forgot password routes
	router.POST("/api/auth/forgot_password", requests.ForgotPassword)
	router.POST("/api/auth/forgot_password/validate_reset_code", requests.ValidateResetCode)
	router.POST("/api/auth/forgot_password/reset_password", requests.ResetPassword)
	router.POST("/api/auth/forgot_password/resend_reset_code", requests.ResendResetCode)

	// events routes
	router.POST("/api/event/add_event", middleware.RequiredAuth, requests.AddEvent)
	router.GET("/api/event/get_all_user_events", middleware.RequiredAuth, requests.GetAllUserEvents)
	router.GET("api/event/get_all_participants_in_event", middleware.RequiredAuth, requests.GetAllParticipatesInEvent)

	// profile routes
	router.POST("api/profile/upload_profile_image", middleware.RequiredAuth, requests.UploadUserProfilePicture)

	return router
}
