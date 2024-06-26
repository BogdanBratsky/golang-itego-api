package main

import (
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/api"
	"github.com/bogdanbratsky/golang-itego-api/config"
	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {
	repositories.InitDB()
	defer repositories.CloseDB()

	if err := config.InitConfig(); err != nil {
		log.Println("Не удалось прочитать config.yaml:", err)
	}

	// Создаем настройки CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{viper.GetString("cors.addres")},                     // Разрешаем запросы только с этого источника
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}, // Разрешаемые HTTP-методы
		AllowedHeaders: []string{"Content-Type", "Authorization"},                    // Разрешенные заголовки
	}

	// Создаем CORS middleware с настройками
	c := cors.New(corsOptions)

	// Применяем CORS middleware к обработчику маршрутов
	handlerWithCors := c.Handler(api.Router())

	// Запускаем HTTP-сервер с обработчиком маршрутов, содержащим CORS middleware
	if err := http.ListenAndServe(viper.GetString("port"), handlerWithCors); err != nil {
		log.Println("Error: ", err)
		return
	}
}
