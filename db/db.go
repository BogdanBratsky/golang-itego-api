package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
)

var DB *sql.DB

func createDSN() string {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		"postgres",
		"2048",
		"localhost",
		"5432",
		"itegodb",
		"disable",
	)
	return dsn
}

func InitDB() {
	var err error
	DB, err = sql.Open("pgx", createDSN())
	if err != nil {
		log.Println("Не удалось подключиться к базе данных: ", err)
		return
	} else {
		log.Println("Connected to DB")
	}

	if err = DB.Ping(); err != nil {
		log.Println("Error connecting to the database:", err)
		return
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Connection was closed")
	}
}
