package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		urlMap     map[string]string
		wantStatus int
		wantHeader string
	}{
		{
			name:    "test 1: valid redirection",
			request: "/id123",
			urlMap: map[string]string{
				"id123": "http://example.com",
			},
			wantStatus: http.StatusTemporaryRedirect,
			wantHeader: "http://example.com",
		},
		{
			name:       "test 2: not found",
			request:    "/nonexistent",
			urlMap:     map[string]string{},
			wantStatus: http.StatusNotFound,
			wantHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем карту urlMap, используем синхронизацию
			urlMap = tt.urlMap

			// Создаем запрос
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			// Обрабатываем запрос
			GetHandler(w, request)

			result := w.Result()

			// Проверяем статус ответа
			assert.Equal(t, tt.wantStatus, result.StatusCode)

			// Проверяем заголовок Location для редиректа
			if tt.wantStatus == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.wantHeader, result.Header.Get("Location"))
			}
		})
	}
}
