package main

import (
    "time"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "sci-abo-go/middleware"
    "sci-abo-go/requests"
)

func InitializerRoutes() *gin.Engine {
    // create a Gin router with default middleware (logger and recovery)
    router := gin.Default()

    // configure CORS middleware
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:4200", "http://localhost:8080"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

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
    router.POST("/api/event/join_event", middleware.RequiredAuth, requests.JoinEvent)
    router.POST("/api/event/upload_event_image", middleware.RequiredAuth, requests.UploadEventPic)
    router.GET("/api/event/search_event", middleware.RequiredAuth, requests.SearchEvent)
    router.GET("/api/event/get_all_user_events", middleware.RequiredAuth, requests.GetAllUserEvents)
    router.GET("api/event/get_all_participants_in_event", middleware.RequiredAuth, requests.GetAllParticipatesInEvent)
    router.GET("api/event/get_event_by_id/", middleware.RequiredAuth, requests.GetEventByID)

    // profile routes
    router.POST("api/profile/upload_profile_image", middleware.RequiredAuth, requests.UploadUserProfilePicture)
    router.POST("api/profile/update_user_details", middleware.RequiredAuth, requests.UpdateUserDetails)
    router.DELETE("api/profile/delete_user", middleware.RequiredAuth, requests.DeleteUser)

    // messages
    // router.POST("api/messages/send_first_message", middleware.RequiredAuth, requests.SendFirstMessage)
    // router.POST("api/messages/send_message", middleware.RequiredAuth, requests.SendAndInsertNewMessageToRoom)
    // router.GET("api/messages/get_all_rooms", middleware.RequiredAuth, requests.GetAllRoomsByUserID)
    // router.GET("api/messages/get_all_messages_by_room_id", middleware.RequiredAuth, requests.GetMessagesByRoomID)

    return router
}
