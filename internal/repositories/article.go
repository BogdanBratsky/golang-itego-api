package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

func ArticleExists(id int) bool {
	query := `
		SELECT COUNT(*)
		FROM article
		WHERE article_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(id).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	return count > 0
}

func AddArticleToDB(a models.Article) error {
	if a.ArticleTitle == "" || a.ArticleText == "" {
		return errors.New("нельзя передавать пустые значения")
	}

	query := `
		INSERT INTO article (article_title, article_text, author_id, category_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return errors.New("не удалось подготовить sql-запрос")
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.ArticleTitle, a.ArticleText, a.AuthorId, a.CategoryId, time.Now())
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return errors.New("не удалось выполнить sql-операцию")
	}

	log.Println("sql-запрос успешно завершён")
	return nil
}

func DeleteArticleFromDB(id int) error {
	query := `
		DELETE FROM article
		WHERE article_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return errors.New("не удалось подготовить sql-запрос")
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return errors.New("не удалось выполнить sql-операцию")
	}

	log.Println("sql-запрос успешно завершён")
	return nil
}

func UpdateArticleInDB(id int, category int, title, text string) error {
	query := `
		UPDATE article 
		SET article_title = $1,
			article_text = $2,
			category_id = $3
		WHERE article_id = $4
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return errors.New("не удалось подготовить sql-запрос")
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, text, category, id)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return errors.New("не удалось выполнить sql-операцию")
	}

	log.Println("sql-запрос успешно завершён")
	return nil
}

func GetArticlesFromDB(page, limit int) ([]models.Article, int, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	if limit >= 30 {
		limit = 30
	}
	offset := (page - 1) * limit

	query := `
		SELECT *
		FROM article
		ORDER BY article_id DESC
		LIMIT $1 OFFSET $2
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		err := rows.Scan(
			&a.ArticleId,
			&a.ArticleTitle,
			&a.ArticleText,
			&a.ArticleViews,
			&a.AuthorId,
			&a.CategoryId,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, a)
	}

	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM article`).Scan(&count)

	return articles, count, nil
}

func GetArticleFromDB(id int) (models.Article, error) {
	query := `
		SELECT * 
		FROM article
		WHERE article_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return models.Article{}, err
	}
	defer stmt.Close()

	var a models.Article
	err = stmt.QueryRow(id).Scan(
		&a.ArticleId,
		&a.ArticleTitle,
		&a.ArticleText,
		&a.ArticleViews,
		&a.AuthorId,
		&a.CategoryId,
		&a.CreatedAt,
	)
	if err != nil {
		return models.Article{}, err
	}

	return a, nil
}

func GetArticleByCategoryFromDB(id, page, limit int) ([]models.Article, int, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	if limit >= 30 {
		limit = 30
	}
	offset := (page - 1) * limit

	query := `
		SELECT * 
		FROM article
		WHERE category_id = $1
		ORDER BY article_id DESC
		LIMIT $2 OFFSET $3
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		err := rows.Scan(
			&a.ArticleId,
			&a.ArticleTitle,
			&a.ArticleText,
			&a.ArticleViews,
			&a.AuthorId,
			&a.CategoryId,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, a)
	}

	log.Println(articles)

	var count int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM article WHERE category_id = $1`, id).Scan(&count); err != nil {
		log.Println(err)
		return nil, 0, err
	}

	return articles, count, nil
}
