package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "test 1: successful URL shortening",
			method:         http.MethodPost,
			body:           "http://example.com",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "test 2: empty URL",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "URL is empty\n",
		},
		{
			name:           "test 3: invalid method",
			method:         http.MethodGet,
			body:           "http://example.com",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Invalid request method\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			// Обрабатываем запрос
			PostHandler(w, request)

			result := w.Result()

			// Проверяем статус ответа
			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			// Проверяем тело ответа только для тех случаев, где оно ожидается
			if tt.expectedBody != "" {
				body, _ := io.ReadAll(result.Body)
				assert.Equal(t, tt.expectedBody, string(body))
			}

			// Проверяем, что сокращенный URL возвращается в случае успешного запроса
			if tt.expectedStatus == http.StatusCreated {
				body, _ := io.ReadAll(result.Body)
				assert.Contains(t, string(body), "http://localhost:8080/")
			}
		})
	}
}
