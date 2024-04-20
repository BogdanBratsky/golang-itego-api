package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bogdanbratsky/golang-itego-api/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/spf13/viper"
)

var DB *sql.DB

type dbConfig struct {
	user     string
	password string
	host     string
	port     string
	dbname   string
	sslmode  string
}

func initDSN() string {
	if err := config.InitConfig(); err != nil {
		log.Println("Не удалось прочитать config.yaml:", err)
	}

	dbConfig := dbConfig{
		user:     viper.GetString("database.user"),
		password: viper.GetString("database.password"),
		host:     viper.GetString("database.host"),
		port:     viper.GetString("database.port"),
		dbname:   viper.GetString("database.dbname"),
		sslmode:  viper.GetString("database.sslmode"),
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		dbConfig.user,
		dbConfig.password,
		dbConfig.host,
		dbConfig.port,
		dbConfig.dbname,
		dbConfig.sslmode,
	)
}

func InitDB() {
	var err error
	DB, err = sql.Open("pgx", initDSN())
	if err != nil {
		log.Println("Не удалось соединиться с базой данных:", err)
		return
	} else {
		log.Println("Соединение с базой данных установлено")
	}

	if err = DB.Ping(); err != nil {
		log.Println("Ошибка соединения с базой данных:", err)
		return
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Соединение с базой данных было закрыто")
	}
}
