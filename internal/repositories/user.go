package repositories

import (
	"errors"
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
func UserExists(name, password string) (int, bool) {
	log.Printf("Проверка пользователя с именем: %s и паролем: %s\n", name, password)

	query := `
		SELECT member_id, superuser
		FROM member 
		WHERE 
			member_name = $1
			AND member_password = $2
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return 0, false
	}
	defer stmt.Close()

	var userId int
	var superuser bool
	err = stmt.QueryRow(name, password).Scan(&userId, &superuser)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return 0, false
	} else {
		log.Println("sql-операция выполнена успешно")
		return userId, superuser
	}
}

// функция для проверки существования пользователя
func UserExistsByID(id int) bool {
	query := `
		SELECT COUNT(*)
		FROM member 
		WHERE member_id = $1
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
	}

	log.Println("sql-операция выполнена успешно")
	return id
}

// функция для получения из базы данных среза всех пользователей и их общего количества
func GetUsersFromDB() ([]models.UserDTO, int) {
	query := `
		SELECT member_id, member_name, member_email, superuser 
		FROM member 
		WHERE superuser != true
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

func GetUserFromDB(id int) models.UserDTO {
	query := `
		SELECT member_id, member_name, member_email, superuser
		FROM member
		WHERE member_id = $1
	`
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("не удалось подготовить sql-запрос:", err)
		return models.UserDTO{}
	}
	defer stmt.Close()

	var user models.UserDTO
	err = stmt.QueryRow(id).Scan(
		&user.UserId,
		&user.UserName,
		&user.UserEmail,
		&user.Superuser,
	)
	if err != nil {
		log.Println("не удалось выполнить sql-операцию:", err)
		return models.UserDTO{}
	}

	return user
}

func DeleteUserFromDB(id int) error {
	query := `
		DELETE FROM member
		WHERE member_id = $1
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
