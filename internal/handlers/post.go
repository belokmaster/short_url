package handlers

import (
	"io"        // Для чтения тела запроса
	"math/rand" // Для генерации случайных чисел
	"net/http"  // Для создания HTTP-сервера
	"sync"
)

// Карту для хранения соответствий сокращённых идентификаторов и оригинальных URL
var urlMap = make(map[string]string)
var mu sync.Mutex

// Функция generateShortID генерирует случайный идентификатор заданной длины.
func generateShortID(lenght int) (string, error) {
	// charset - набор символов, которые могут быть использованы в ID.
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Создаём срез байтов нужной длины.
	b := make([]byte, lenght)

	// Чтение случайных байтов из криптографически безопасного источника.
	if _, err := rand.Read(b); err != nil {
		// Если не удалось прочитать случайные байты, возвращаем ошибку.
		return "", err
	}

	// Проходим по всем сгенерированным байтам и маппируем их на символы из charset.
	for i := range b {
		// Используем остаток от деления случайного числа на длину charset для получения символа.
		b[i] = charset[rand.Intn(len(charset))]
	}
	// Преобразуем срез байтов в строку и возвращаем.
	return string(b), nil
}

// Функция shortenURL генерирует сокращённый URL с использованием заданного идентификатора.
func shortenURL(url string) (string, error) {
	// Генерация случайной длины для идентификатора.
	length := rand.Intn(len(url)) + 4 // rand.Intn(8) генерирует число от 0 до длины url, добавляем 4 по приколу

	shortID, err := generateShortID(length)
	if err != nil {
		return "", err
	}

	// Сохраняем соответствие между shortID и оригинальным URL в карту
	mu.Lock()
	urlMap[shortID] = url
	mu.Unlock()

	return "http://localhost:8080/" + shortID, nil
}

// Функция PostHandler обрабатывает POST-запросы на корневом пути ("/").
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка, что метод запроса - POST. Если нет, возвращается ошибка 405 (Method Not Allowed).
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Чтение тела запроса
	// io.ReadAll читает все данные из r.Body и возвращает их в виде среза байтов.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Если возникла ошибка при чтении тела запроса, возвращается ошибка 400 (Bad Request).
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Закрытие тела запроса для освобождения ресурсов.
	defer r.Body.Close()

	// Преобразование тела запроса из среза байтов в строку.
	url := string(body)

	// Проверка, что URL не пустой.
	if url == "" {
		// Если URL пустой, возвращается ошибка 400 (Bad Request).
		http.Error(w, "URL is empty", http.StatusBadRequest)
		return
	}

	// Сокращение URL с помощью функции shortenURL.
	shortenedURL, err := shortenURL(url)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	// Установка заголовка Content-Type в "text/plain".
	w.Header().Set("Content-Type", "text/plain")
	// Установка статуса ответа в 201 (Created).
	w.WriteHeader(http.StatusCreated)

	// Запись сокращённого URL в тело ответа.
	// w.Write принимает срез байтов, поэтому строку необходимо преобразовать в срез байтов.
	w.Write([]byte(shortenedURL))
}
