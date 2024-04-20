package main

import (
	"log"
	"net/http"

	"github.com/bogdanbratsky/golang-itego-api/api"
	"github.com/bogdanbratsky/golang-itego-api/config"
	"github.com/bogdanbratsky/golang-itego-api/internal/repositories"
	"github.com/spf13/viper"
)

func main() {
	repositories.InitDB()
	defer repositories.CloseDB()

	if err := config.InitConfig(); err != nil {
		log.Println("Не удалось прочитать config.yaml:", err)
	}

	if err := http.ListenAndServe(viper.GetString("port"), api.Router()); err != nil {
		log.Println("Error: ", err)
		return
	}
}
