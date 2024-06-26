package repositories

import (
	"errors"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

func CategoryExists(id int) bool {
	query := `
		SELECT COUNT(*)
		FROM category
		WHERE category_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(id).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func AddCategoryToDB(c models.Category) error {
	query := `
		INSERT INTO category (category_name)
		VALUES ($1)
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.CategoryName)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCategoryFromDB(id int) error {
	query := `
		DELETE FROM category WHERE category_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCategoryInDB(id int, c models.Category) error {
	if c.CategoryName == "" {
		return errors.New("нельзя передавать пустые значения")
	}

	query := `
		UPDATE category SET category_name = $1 WHERE category_id = $2
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.CategoryName, id)
	if err != nil {
		return err
	}

	return nil
}

func GetCategoriesFromDB() ([]models.Category, int, error) {
	query := `
		SELECT *
		FROM category
		ORDER BY category_id DESC
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.CategoryId,
			&c.CategoryName,
		)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}

	var count int
	err = DB.QueryRow(`SELECT COUNT(*) FROM category`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func GetCategoryFromDB(id int) (models.Category, error) {
	query := `
		SELECT *
		FROM category
		WHERE category_id = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return models.Category{}, err
	}
	defer stmt.Close()

	var c models.Category
	err = stmt.QueryRow(id).Scan(&c.CategoryId, &c.CategoryName)
	if err != nil {
		return models.Category{}, err
	}

	return c, nil
}
