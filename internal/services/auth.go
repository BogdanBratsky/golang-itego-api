package services

import (
	"errors"

	"github.com/bogdanbratsky/golang-itego-api/models"
)

func CreateUser(name, email, password string) (*models.UserDTO, error) {
	if name == "" || email == "" || password == "" {
		return &models.UserDTO{}, errors.New("нельзя передавать пустые значения")
	}

	return &models.UserDTO{
		UserId:    1,
		UserName:  name,
		UserEmail: email,
	}, nil
}
