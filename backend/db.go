package main

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// Конфигурация подключения к БД
const (
	host     = "localhost"   // Хост БД
	port     = 5432          // Порт PostgreSQL
	user     = "postgres"      // Имя пользователя БД
	password = "postgres"  // Пароль пользователя
	dbname   = "devops_app"  // Имя базы данных
)

// Функция подключения к БД
func connectDB() *sql.DB {
	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	// Открытие соединения
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	
	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Successfully connected to DB!")
	return db
}
