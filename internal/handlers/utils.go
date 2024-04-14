package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

// Структура для формирования успешного ответа
type response struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Status  int              `json:"status"`
	User    *models.UserDTO  `json:"user,omitempty"`
	Token   string           `json:"token,omitempty"`
	Users   []models.UserDTO `json:"users,omitempty"`
}

// Функция для формирования успешного ответа
func RespondWithSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Не удалось сериализовать ответ: ", err)
		return
	}
}

// Функция для формирования ответа с ошибкой
func RespondWithError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Success: false,
		Message: message,
		Status:  status,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Не удалось сериализовать ответ: ", err)
		return
	}
}
