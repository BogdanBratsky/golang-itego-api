package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/bogdanbratsky/golang-itego-api/models"
)

func CreateUser(u models.User) (models.UserDTO, error) {
	// проверяем, чтобы строки не были пустыми
	if u.UserName == "" || u.UserEmail == "" || u.UserPassword == "" {
		return models.UserDTO{}, errors.New("нельзя передавать пустые значения")
	}
	// проверяем, чтоб длина пароля была не меньше 8 символов
	if len([]byte(u.UserPassword)) < 8 {
		return models.UserDTO{}, errors.New("пароль должен содержать минимум 8 символов")
	}
	// проверяем, занято ли имя
	if !repositories.IsNameTaken(u.UserName) {
		return models.UserDTO{}, errors.New("имя занято")
	}
	// проверяем, занята ли почта
	if !repositories.IsEmailTaken(u.UserEmail) {
		return models.UserDTO{}, errors.New("почта занята")
	}

	// хэшируем пароль
	hashedPassword := HashPassword(u.UserPassword)

	// после успешно завершённых проверок и хэширования пароля отдаём данные для записи в бд
	u.UserId = repositories.AddUserToDB(
		u.UserName,
		u.UserEmail,
		hashedPassword,
		u.Superuser,
	)

	// возвращаем структуру с информацией о созданном пользователе для ответа
	return models.UserDTO{
		UserId:    u.UserId,
		UserName:  u.UserName,
		UserEmail: u.UserEmail,
		Superuser: u.Superuser,
	}, nil
}

// Map для хранения времени последней попытки входа и количества попыток с каждого IP
var loginAttempts = make(map[string]*AttemptInfo)
var mu sync.Mutex

const (
	maxAttempts     = 5             // Максимальное количество попыток
	attemptInterval = 1 * time.Hour // Время, в течение которого учитываются попытки
)

type AttemptInfo struct {
	LastAttempt time.Time
	Count       int
}

func SignInUser(u models.User, ip string) (*models.UserDTO, error) {
	if u.UserName == "" || u.UserPassword == "" {
		return &models.UserDTO{}, errors.New("нельзя передавать пустые значения")
	}

	mu.Lock()
	attemptInfo, exists := loginAttempts[ip]
	if !exists {
		attemptInfo = &AttemptInfo{LastAttempt: time.Now(), Count: 0}
		loginAttempts[ip] = attemptInfo
	}
	// Обновление информации о попытках
	if time.Since(attemptInfo.LastAttempt) > attemptInterval {
		attemptInfo.Count = 0
	}
	attemptInfo.LastAttempt = time.Now()
	attemptInfo.Count++
	mu.Unlock()

	// Проверка на превышение количества попыток
	if attemptInfo.Count > maxAttempts {
		return &models.UserDTO{}, errors.New("слишком много попыток входа, попробуйте позже")
	}

	// Хэшируем пароль для того, чтобы сверить его с паролем, записанным в бд в виде хэша
	hashedPassword := HashPassword(u.UserPassword)

	userId, superuser := repositories.UserExists(
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
		Superuser: superuser,
	}, nil
}

// HashPassword хэширует переданный пароль и возвращает его хэш
func HashPassword(password string) string {
	hashed := sha256.Sum256([]byte(password + "dzLKCTDJVDaR"))
	hashedString := fmt.Sprintf("%x", hashed)
	return hashedString
}
