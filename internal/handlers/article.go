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

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
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

	var article models.Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "не удалось декодировать json", http.StatusBadRequest)
		return
	}

	err = services.CreateArticle(article, token)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "не удалось внести запись в бд", http.StatusInternalServerError)
		return
	}

	response := response{
		Success: true,
		Message: "запись успешно создана",
		Status:  http.StatusCreated,
	}

	RespondWithSuccess(w, http.StatusCreated, response)
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	articleId, _ := strconv.Atoi(vars["id"])

	if articleId <= 0 {
		log.Println("id не может быть меньше либо равно 0")
		RespondWithError(w, "id не может быть меньше либо равно 0", http.StatusBadRequest)
		return
	}

	if !repositories.ArticleExists(articleId) {
		log.Println("записи не существует")
		RespondWithError(w, "записи не существует", http.StatusBadRequest)
		return
	}

	err = repositories.DeleteArticleFromDB(articleId)
	if err != nil {
		log.Println(err)
		RespondWithError(w, "не удалось удалить запись из бд", http.StatusInternalServerError)
		return
	}

	response := response{
		Success: true,
		Message: "запись успешно удалена",
		Status:  http.StatusNoContent,
	}

	RespondWithSuccess(w, http.StatusCreated, response)
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
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
	articleId, _ := strconv.Atoi(vars["id"])

	if !repositories.ArticleExists(articleId) {
		log.Println("записи не существует")
		RespondWithError(w, "записи не существует", http.StatusBadRequest)
		return
	}

	var article models.Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println(err)
		return
	}

	err = repositories.UpdateArticleInDB(
		articleId,
		article.CategoryId,
		article.ArticleTitle,
		article.ArticleText,
	)
	if err != nil {
		log.Println(err)
		return
	}

	response := response{
		Success: true,
		Message: "запись успешно изменена",
		Status:  http.StatusOK,
	}

	RespondWithSuccess(w, http.StatusCreated, response)
}

func GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	articles, count, err := repositories.GetArticlesFromDB(page, limit)
	if err != nil {
		log.Println(err)
		RespondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := response{
		Success:    true,
		Message:    "список записей успешно получен",
		Status:     http.StatusOK,
		Articles:   articles,
		TotalCount: count,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}

func GetArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	vars := mux.Vars(r)
	articleId, _ := strconv.Atoi(vars["id"])

	article, err := repositories.GetArticleFromDB(articleId)
	if err != nil {
		log.Println(err)
		return
	}

	response := response{
		Success: true,
		Message: "запись успешно получена",
		Status:  http.StatusOK,
		Article: &article,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}

func GetArticlesByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request recieved:", r.Method, r.URL)

	vars := mux.Vars(r)
	categoryId, _ := strconv.Atoi(vars["id"])

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	articles, count, err := repositories.GetArticleByCategoryFromDB(
		categoryId,
		page,
		limit,
	)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(articles)

	response := response{
		Success:    true,
		Message:    "запись успешно получена",
		Status:     http.StatusOK,
		Articles:   articles,
		TotalCount: count,
	}

	RespondWithSuccess(w, http.StatusOK, response)
}
