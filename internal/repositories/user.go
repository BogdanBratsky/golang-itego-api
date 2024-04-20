package repositories

import (
	"log"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

// функция для проверки не занято ли имя пользователя
func IsNameTaken(name string) bool {
	query := `
		SELECT COUNT(*)
		FROM member
		WHERE member_name = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(name).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

// функция для проверки не занята ли почта
func IsEmailTaken(email string) bool {
	query := `
		SELECT COUNT(*)
		FROM member
		WHERE member_email = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(email).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

// функция для проверки не занят ли пароль
func IsPasswordTaken(password string) bool {
	query := `
		SELECT COUNT(*)
		FROM member
		WHERE member_password = $1
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(password).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

// функция для проверки существования пользователя
func UserExists(name, password string) int {
	query := `
		SELECT member_id 
		FROM member 
		WHERE 
			member_name = $1
			AND member_password = $2
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return 0
	}
	defer stmt.Close()

	var userId int
	err = stmt.QueryRow(name, password).Scan(&userId)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return 0
	} else {
		log.Println("sql-операция выполнена успешно")
		return userId
	}
}

// функция для добавления нового пользователя в базу данных
func AddUserToDB(name, email, password string, superuser bool) int {
	query := `
		INSERT INTO member
			(member_name, member_email, member_password, superuser)
		VALUES
			($1, $2, $3, $4)
		RETURNING member_id
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return 0
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, email, password, superuser).Scan(&id)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return 0
	} else {
		log.Println("sql-операция выполнена успешно")
		return id
	}
}

// функция для получения из базы данных среза всех пользователей и их общего количества
func GetUsersFromDB() ([]models.UserDTO, int) {
	query := `
		SELECT member_id, member_name, member_email, superuser 
		FROM member 
		ORDER BY member_id DESC
	`
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return nil, 0
	}
	defer rows.Close()

	var users []models.UserDTO
	for rows.Next() {
		var u models.UserDTO
		err := rows.Scan(&u.UserId, &u.UserName, &u.UserEmail, &u.Superuser)
		if err != nil {
			return nil, 0
		}
		users = append(users, u)
	}
	var count int
	err = DB.QueryRow(`SELECT COUNT(*) FROM member`).Scan(&count)
	if err != nil {
		return nil, 0
	}

	return users, count
}

// func GetUserFromDB() models.UserDTO {
// 	query := `
// 	SELECT user_id, user_name, user_email, superuser
// 	FROM users
// 	WHERE user_id = $1
// 	ORDER BY user_id DESC
// 	`
// 	rows, err := DB.QueryRow(query)
// 	if err != nil {
// 		log.Println("не удалось выполнить sql-операцию:", err)
// 		return models.UserDTO{}
// 	}
// 	defer rows.Close()

// 	var users []models.UserDTO
// 	for rows.Next() {
// 		var u models.UserDTO
// 		err := rows.Scan(&u.UserId, &u.UserName, &u.UserEmail, &u.Superuser)
// 		if err != nil {
// 			return nil
// 		}
// 		users = append(users, u)
// 	}
// 	return users
// }
