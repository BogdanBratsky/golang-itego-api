package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

// Структура для формирования успешного ответа
type response struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Status     int               `json:"status"`
	Name       string            `json:"name,omitempty"`
	Superuser  bool              `json:"superuser,omitempty"`
	User       *models.UserDTO   `json:"user,omitempty"`
	Token      string            `json:"token,omitempty"`
	Users      []models.UserDTO  `json:"users,omitempty"`
	Article    *models.Article   `json:"article,omitempty"`
	Articles   []models.Article  `json:"articles,omitempty"`
	Category   *models.Category  `json:"category,omitempty"`
	Categories []models.Category `json:"categories,omitempty"`
	TotalCount int               `json:"totalCount,omitempty"`
}

// Функция для формирования успешного ответа
func RespondWithSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Не удалось сериализовать ответ:", err)
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
		log.Println("Не удалось сериализовать ответ:", err)
		return
	}
}

// функция для извлечения токена из запроса
func GetToken(r *http.Request) (string, error) {
	// Получаем значение заголовка Authorization
	authHeader := r.Header.Get("Authorization")

	// Проверяем, что заголовок существует и начинается с "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("не удалось извлечь токен")
	}

	// Извлекаем токен из строки "Bearer <token>"
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
