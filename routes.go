package main

import (
    "sci-abo-go/requests"
    "github.com/gin-gonic/gin"
)

func InitializerRoutes() *gin.Engine {
    // Create a Gin router with default middleware (logger and recovery)
    router := gin.Default()

    // Auth routes
    router.POST("/auth/register", requests.CreateUser)
    router.POST("/auth/login", requests.Login)
    router.POST("/auth/validate", requests.Validate)

    // Profile routes
    router.POST("/profile/upload_profile_image", requests.UploadUserProfilePicture)

    return router
}
