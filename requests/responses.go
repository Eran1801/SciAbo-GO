package requests

import (
    "encoding/json"
    "net/http"
)

type errorResponse struct {
    Error  string `json:"error"`
}

type successResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"` // omitempty means not required
}

func ErrorResponse(message string, w http.ResponseWriter) {
    response := errorResponse{
        Error:  message,
    }
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(response)
}

func SuccessResponse(message string, data interface{}, w http.ResponseWriter) {
    response := successResponse{
        Message: message,
        Data:    data,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
