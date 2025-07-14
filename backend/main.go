package main

import (
	"database/sql"
	"log"
	"net/http"
)

// Глобальная переменная подключения к БД
var db *sql.DB

func main() {
	// Подключение к БД
	db = connectDB()
	defer db.Close()

	// Обслуживание статических файлов фронтенда
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Обработчик для /items (GET и POST)
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		// Настройка CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Обработка предварительного запроса OPTIONS
		if r.Method == "OPTIONS" {
			return
		}
		
		// Роутинг методов
		switch r.Method {
		case "GET":
			getItems(w, r)
		case "POST":
			createItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для /update-item (PUT)
	http.HandleFunc("/update-item", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			return
		}
		
		if r.Method == "PUT" {
			updateItem(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для /delete-item (DELETE)
	http.HandleFunc("/delete-item", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			return
		}
		
		if r.Method == "DELETE" {
			deleteItem(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Запуск сервера на порту 8080
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
