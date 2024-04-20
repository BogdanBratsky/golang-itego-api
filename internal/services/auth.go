package services

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/bogdanbratsky/golang-itego-api/models"
)

func CreateUser(u models.User) (*models.UserDTO, error) {
	if u.UserName == "" || u.UserEmail == "" || u.UserPassword == "" {
		return &models.UserDTO{}, errors.New("нельзя передавать пустые значения")
	}

	if !repositories.IsNameTaken(u.UserName) {
		return &models.UserDTO{}, errors.New("с таким именем уже зарегистрирована учётная запись")
	}
	if !repositories.IsEmailTaken(u.UserEmail) {
		return &models.UserDTO{}, errors.New("на эту почту уже зарегистрирована учётная запись")
	}
	if !repositories.IsPasswordTaken(u.UserPassword) {
		return &models.UserDTO{}, errors.New("с таким паролем уже зарегистрирована учётная запись")
	}

	// хэшируем пароль
	hashedPassword := HashPassword(u.UserPassword)

	// отправляем данные в бд
	u.UserId = repositories.AddUserToDB(
		u.UserName,
		u.UserEmail,
		hashedPassword,
		u.Superuser,
	)

	// возвращаем структуру с информацией о созданном пользователе для ответа
	return &models.UserDTO{
		UserId:    u.UserId,
		UserName:  u.UserName,
		UserEmail: u.UserEmail,
		Superuser: u.Superuser,
	}, nil
}

func SignInUser(u models.User) (*models.UserDTO, error) {
	if u.UserName == "" || u.UserPassword == "" {
		return &models.UserDTO{}, errors.New("нельзя передавать пустые значения")
	}

	// хэшируем пароль для того, чтобы сверить его с паролем, записанным в бд в виде хэша
	hashedPassword := HashPassword(u.UserPassword)

	userId := repositories.UserExists(
		u.UserName,
		hashedPassword,
	)

	if userId == 0 {
		return &models.UserDTO{}, errors.New("пользователь не существует")
	}

	return &models.UserDTO{
		UserId:    userId,
		UserName:  u.UserName,
		UserEmail: u.UserEmail,
		Superuser: u.Superuser,
	}, nil
}

// HashPassword хэширует переданный пароль и возвращает его хэш
func HashPassword(password string) string {
	hashed := sha256.Sum256([]byte(password + "salt"))
	hashedString := fmt.Sprintf("%x", hashed)
	return hashedString
}
