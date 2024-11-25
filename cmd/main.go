package main

import (
	"fmt"                         // Для чтения тела запроса
	"net/http"                    // Для создания HTTP-сервера
	"short_url/internal/handlers" // Импортируем пакет handlers из папки internal
)

func main() {
	// Регистрация обработчиков
	http.HandleFunc("/", handlers.PostHandler)        // для POST-запросов на корневом пути
	http.HandleFunc("/shorten", handlers.PostHandler) // для POST-запросов на пути "/shorten"
	http.HandleFunc("/", handlers.GetHandler)         // для GET-запросов на пути "/{id}"

	// Вывод сообщения в консоль о запуске сервера.
	fmt.Println("Server is running on port 8080")
	// Если возникает ошибка при запуске сервера, выводится сообщение об ошибке.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
