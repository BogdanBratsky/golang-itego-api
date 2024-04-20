package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/bogdanbratsky/golang-itego-api/internal/services"
	"github.com/bogdanbratsky/golang-itego-api/models"
)

// Обработчик для регистрации пользователя
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request recieved:", r.Method, r.URL)

	// Получаем данные для регистрации от пользователя
	// в формате JSON и десериализуем их в структуру User
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("не удалось декодировать тело запроса:", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Полученные и преобразованные данные отдаём на слой бизнес-логики (services)
	newUser, err := services.CreateUser(user)
	if err != nil {
		log.Println("не удалось создать пользователя:", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "пользователь успешно создан",
		Status:  http.StatusCreated,
		User:    newUser,
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusCreated, response)
}

// Обработчик для авторизации пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("не удалось декодировать тело запроса:", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDTO, err := services.SignInUser(user)
	if err != nil {
		log.Println("не удалось авторизоваться:", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "пользователь успешно авторизован",
		Status:  http.StatusOK,
		User:    userDTO,
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusOK, response)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	users, count := repositories.GetUsersFromDB()

	// Структура для формирования успешного ответа
	response := response{
		Success:    true,
		Message:    "список пользователей успешно получен",
		Status:     http.StatusOK,
		Users:      users,
		TotalCount: count,
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusOK, response)
}
