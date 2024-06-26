package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/bogdanbratsky/golang-itego-api/internal/services"
	"github.com/bogdanbratsky/golang-itego-api/models"
	"github.com/gorilla/mux"
)

// Обработчик для регистрации пользователя
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request recieved:", r.Method, r.URL)

	token, err := GetToken(r)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	success, err := services.ValidateToken(token)
	if err != nil {
		log.Println(err)
		return
	}
	if !success {
		log.Println("токен недействителен")
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	payload, err := services.ExtractPayload(token)
	if err != nil {
		log.Println(err)
		return
	}
	isSuperuser := payload["superuser"].(bool)
	if !isSuperuser {
		log.Println("нет прав на это действие")
		RespondWithError(w, "нет прав на это действие", http.StatusForbidden)
		return
	}

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
	log.Println("создана новая учётная запись:", newUser)

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "пользователь успешно создан",
		Status:  http.StatusCreated,
		// User:    newUser,
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusCreated, response)
}

// Обработчик для авторизации пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received:", r.Method, r.URL)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("не удалось декодировать тело запроса:", err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr

	userDTO, err := services.SignInUser(user, ip)
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
		Name:    userDTO.UserName,
		// User:    userDTO,
		Token: services.GenerateToken(
			userDTO.UserId,
			userDTO.UserName,
			userDTO.Superuser,
		),
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusOK, response)
}

// Обработчик для проверки действительности токена
func CheckToken(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	token, err := GetToken(r)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusForbidden)
		return
	}

	success, err := services.ValidateToken(token)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusForbidden)
		return
	}
	if !success {
		log.Println("токен недействителен")
		RespondWithError(w, "токен недействителен", http.StatusForbidden)
		return
	}

	payload, err := services.ExtractPayload(token)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(payload)
	isSuperuser := payload["superuser"].(bool)
	if !isSuperuser {
		// Структура для формирования успешного ответа
		resp := struct {
			Success   bool   `json:"success"`
			Message   string `json:"message"`
			Status    int    `json:"status"`
			Superuser bool   `json:"superuser"`
		}{
			Success:   true,
			Message:   "токен действителен",
			Status:    http.StatusOK,
			Superuser: isSuperuser,
		}

		// Формируем успешный ответ
		RespondWithSuccess(w, http.StatusOK, resp)
		return
	}

	// Структура для формирования успешного ответа
	response := response{
		Success:   true,
		Message:   "токен действителен",
		Status:    http.StatusOK,
		Superuser: isSuperuser,
		// User:    &user,
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

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	userId, _ := strconv.Atoi(userIdStr)

	user := repositories.GetUserFromDB(userId)

	// Структура для формирования успешного ответа
	response := response{
		Success: true,
		Message: "пользователь успешно получен",
		Status:  http.StatusOK,
		User:    &user,
	}

	// Формируем успешный ответ
	RespondWithSuccess(w, http.StatusOK, response)
}

// Обработчик для удаления пользователя
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	token, err := GetToken(r)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	success, err := services.ValidateToken(token)
	if err != nil {
		log.Println(err)
		return
	}
	if !success {
		log.Println("токен недействителен")
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	payload, err := services.ExtractPayload(token)
	if err != nil {
		log.Println(err)
		return
	}
	isSuperuser := payload["superuser"].(bool)
	if !isSuperuser {
		log.Println("нет прав на это действие")
		RespondWithError(w, "нет прав на это действие", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])

	if userId <= 0 {
		log.Println("id не может быть меньше либо равно 0")
		RespondWithError(w, "id не может быть меньше либо равно 0", http.StatusBadRequest)
		return
	}

	if !repositories.UserExistsByID(userId) {
		log.Println("записи не существует")
		RespondWithError(w, "записи не существует", http.StatusBadRequest)
		return
	}

	err = repositories.DeleteUserFromDB(userId)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "не удалось удалить запись из бд", http.StatusInternalServerError)
		return
	}

	response := response{
		Success: true,
		Message: "пользователь успешно удалён",
		Status:  http.StatusNoContent,
	}

	RespondWithSuccess(w, http.StatusCreated, response)
}
