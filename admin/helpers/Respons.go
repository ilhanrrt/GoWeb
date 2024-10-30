package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"` // Başarılı yanıtta veri döndürmek için
	Meta    Meta        `json:"meta"`
}

type Meta struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	ErrorMessage   string `json:"errorMessage,omitempty"` // Hata varsa bu mesaj
	ErrorCode      string `json:"errorCode,omitempty"`    // Hata kodu varsa
	Message        string `json:"message"`
}

// Başarılı yanıt
func SuccessResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	response := Response{
		Success: true,
		Data:    data,
		Meta: Meta{
			Message:        message,
			HTTPStatusCode: statusCode,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Hatalı yanıt
func ErrorResponse(w http.ResponseWriter, errorMessage string, errorCode string, statusCode int) {
	response := Response{
		Success: false,
		Meta: Meta{
			HTTPStatusCode: statusCode,
			ErrorMessage:   errorMessage,
			ErrorCode:      errorCode,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
