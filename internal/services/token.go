package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("yek-terces-ogeti")

// TokenGenerator генерирует токен
func GenerateToken(id int, name string, superuser bool) string {
	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	// Задаем набор клеймов (payload) для токена
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["name"] = name
	claims["superuser"] = superuser
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix() // Устанавливаем срок действия токена на 72 часа

	// Подписываем токен с помощью секретного ключа
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("ошибка при подписи токена:", err)
		return ""
	}

	// Возвращаем сгенерированный токен
	return tokenString
}

func ValidateToken(tokenString string) (bool, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return false, nil
	}

	// Проверяем срок действия токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("unable to parse claims")
	}
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if expirationTime.Before(time.Now()) {
		return false, nil
	}

	// Токен валиден
	return true, nil
}

// ExtractPayload извлекает полезную нагрузку (payload) из JWT токена
func ExtractPayload(tokenString string) (map[string]interface{}, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге токена: %v", err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, fmt.Errorf("невалидный токен")
	}

	// Извлекаем полезную нагрузку (payload)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("невозможно извлечь полезную нагрузку")
	}

	// Возвращаем полезную нагрузку
	return claims, nil
}
