package services

import (
	"errors"
	"log"

	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/bogdanbratsky/golang-itego-api/models"
)

func CreateArticle(a models.Article, token string) error {
	if a.ArticleTitle == "" || a.ArticleText == "" {
		return errors.New("нельзя передавать пустые значения")
	}
	// if len([]byte(a.ArticleText)) > 10000 {
	// 	return models.Article{}, errors.New("нельзя передавать пустые значения")
	// }

	payload, _ := ExtractPayload(token)
	id := payload["id"].(float64)

	a.AuthorId = int(id)
	err := repositories.AddArticleToDB(a)
	if err != nil {
		log.Println(err)
		return errors.New("не удалось записать данные в бд")
	}

	return nil
}
