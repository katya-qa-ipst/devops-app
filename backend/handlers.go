package main

import (
	"encoding/json"
	"net/http"
)

// Структура элемента (соответствует таблице в БД)
type Item struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Обработчик GET /items (получение всех элементов)
func getItems(w http.ResponseWriter, r *http.Request) {
	// Запрос к БД
	rows, err := db.Query("SELECT id, title, description FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Сбор результатов
	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Title, &item.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	// Отправка JSON ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Обработчик POST /items (создание элемента)
func createItem(w http.ResponseWriter, r *http.Request) {
	// Парсинг JSON тела запроса
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// SQL запрос на вставку
	sqlStatement := `
	INSERT INTO items (title, description)
	VALUES ($1, $2)
	RETURNING id`
	
	id := 0
	err = db.QueryRow(sqlStatement, newItem.Title, newItem.Description).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Установка ID и отправка ответа
	newItem.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}

// Обработчик PUT /update-item (обновление элемента)
func updateItem(w http.ResponseWriter, r *http.Request) {
	// Получение ID из query параметра
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Парсинг JSON тела запроса
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// SQL запрос на обновление
	sqlStatement := `
		UPDATE items
		SET title = $1, description = $2
		WHERE id = $3`
	res, err := db.Exec(sqlStatement, item.Title, item.Description, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка количества обновленных строк
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if count == 0 {
		http.Error(w, "No item found with given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик DELETE /delete-item (удаление элемента)
func deleteItem(w http.ResponseWriter, r *http.Request) {
	// Получение ID из query параметра
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// SQL запрос на удаление
	sqlStatement := `DELETE FROM items WHERE id = $1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка количества удаленных строк
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if count == 0 {
		http.Error(w, "No item found with given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
