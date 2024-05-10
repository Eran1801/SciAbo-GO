package requests

import (
    "net/http"
    "github.com/gin-gonic/gin"
    
)

func ErrorResponse(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": message,
    })
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "message": message,
        "data": data,
    })
}
