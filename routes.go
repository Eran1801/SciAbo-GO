package main

import (
	"sci-abo-go/middleware"
	"sci-abo-go/requests"

	"github.com/gin-gonic/gin"
)

func InitializerRoutes() *gin.Engine {
    // Create a Gin router with default middleware (logger and recovery)
    router := gin.Default()

    // Auth routes
    router.POST("/auth/register", requests.CreateUser)
    router.POST("/auth/login", requests.Login)
    router.GET("/get_all_events", middleware.RequiredAuth, requests.GetAllEvents)

    // Profile routes
    router.POST("/profile/upload_profile_image", middleware.RequiredAuth, requests.UploadUserProfilePicture)

    return router
}
