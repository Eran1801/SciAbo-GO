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
	router.POST("/auth/register", requests.CreateUser)
	router.POST("/auth/login", requests.Login)
	router.POST("/auth/forget_password", requests.ForgotPassword)
	router.POST("/auth/validate_reset_code", requests.ValidateResetCode)
	router.POST("/auth/reset_password", requests.ResetPassword)


	// events routes
	router.GET("/get_all_events", middleware.RequiredAuth, requests.GetAllEvents)

	// profile routes
	router.POST("/profile/upload_profile_image", middleware.RequiredAuth, requests.UploadUserProfilePicture)

	return router
}
