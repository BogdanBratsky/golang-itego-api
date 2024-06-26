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

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
		return
	}

	if category.CategoryName == "" {
		log.Println("нельзя передавать пустые значения")
		RespondWithError(w, "нельзя передавать пустые значения", http.StatusBadRequest)
		return
	}

	err = repositories.AddCategoryToDB(category)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("категория успешно создана")

	response := response{
		Success: true,
		Message: "категория успешно создана",
		Status:  http.StatusCreated,
	}

	RespondWithSuccess(w, http.StatusCreated, response)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	vars := mux.Vars(r)
	categoryId, _ := strconv.Atoi(vars["id"])

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
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	if !repositories.CategoryExists(categoryId) {
		log.Println("записи не существует")
		RespondWithError(w, "категория не существует", http.StatusBadRequest)
		return
	}

	err = repositories.DeleteCategoryFromDB(categoryId)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response{
		Success: true,
		Message: "категория успешно удалена",
		Status:  http.StatusOK,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
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
		RespondWithError(w, "unauthorized", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	categoryId, _ := strconv.Atoi(vars["id"])

	if !repositories.CategoryExists(categoryId) {
		log.Println("записи не существует")
		RespondWithError(w, "категория не существует", http.StatusBadRequest)
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
		return
	}

	err = repositories.UpdateCategoryInDB(categoryId, category)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := response{
		Success: true,
		Message: "категория успешно изменена",
		Status:  http.StatusOK,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	categories, count, err := repositories.GetCategoriesFromDB()
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		RespondWithError(w, "не существует ни одной категории", http.StatusOK)
		return
	}

	response := response{
		Success:    true,
		Message:    "список категорий успешно получен",
		Status:     http.StatusOK,
		Categories: categories,
		TotalCount: count,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}

func GetCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	vars := mux.Vars(r)
	categoryId, _ := strconv.Atoi(vars["id"])

	category, err := repositories.GetCategoryFromDB(categoryId)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response{
		Success:  true,
		Message:  "категория успешно получена",
		Status:   http.StatusOK,
		Category: &category,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}
