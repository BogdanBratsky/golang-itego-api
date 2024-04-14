package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/internal/services"
	"github.com/bogdanbratsky/golang-itego-api/models"
)

// Обработчик для регистрации пользователя
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request recieved: ", r.Method, r.URL)

	// Получаем данные для регистрации от пользователя
	// в формате JSON и десериализуем их в структуру User
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Не удалось декодировать тело запроса: ", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Полученные и преобразованные данные отдаём на слой бизнес-логики (services)
	newUser, err := services.CreateUser(
		user.UserName,
		user.UserEmail,
		user.UserPassword,
	)
	if err != nil {
		log.Println("Не удалось создать пользователя: ", err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "Пользователь успешно создан",
		Status:  http.StatusCreated,
		User:    newUser,
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusCreated, response)
}

// Обработчик для авторизации пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request recieved: ", r.Method, r.URL)

	var user models.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Не удалось декодировать тело запроса: ", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "Пользователь успешно авторизован",
		Status:  http.StatusOK,
		User:    &user,
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusOK, response)
}
