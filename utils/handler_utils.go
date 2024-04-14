package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

// Структура для формирования ответа
type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	User    models.UserDTO `json:"user,omitempty"`
	Token   string         `json:"token,omitempty"`
}

// Функция для формирования успешного ответа
func RespondWithSuccess(w http.ResponseWriter, message string, status int, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		User:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Не удалось сериализовать ответ: ", err)
		return
	}
}

// Функция для формирования ответа с ошибкой
func RespondWithError() {}
