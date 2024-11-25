package main

import (
	"net/http" // Для создания HTTP-сервера
)

// Функция getHandler обрабатывает GET-запросы на пути "/{id}".
func getHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение id из URL пути.
	id := r.URL.Path[1:]

	// Поиск оригинального URL по id в карте.
	mu.Lock()
	originalURL, ok := urlMap[id]
	mu.Unlock()

	if !ok {
		// Если id не найден, возвращаем ошибку 404 (Not Found).
		http.NotFound(w, r)
		return
	}

	// Установка заголовка Location и статуса 307 (Temporary Redirect).
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
