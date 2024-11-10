package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	require.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), totalCount)
}

func TestMainHandlerWhenCityNotSupported(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=unknown", nil) // Неподдерживаемый город
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)      // Код ответа 400
	assert.Equal(t, "wrong city value", responseRecorder.Body.String()) // Сообщение об ошибке
}

func TestMainHandlerWhenCountMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow", nil) // Пропущенный параметр count
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)   // Код ответа 400
	assert.Equal(t, "count missing", responseRecorder.Body.String()) // Сообщение об ошибке
}

func TestMainHandlerWhenWrongCountValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=abc&city=moscow", nil) // Неправильное значение count
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)       // Код ответа 400
	assert.Equal(t, "wrong count value", responseRecorder.Body.String()) // Сообщение об ошибке
}
