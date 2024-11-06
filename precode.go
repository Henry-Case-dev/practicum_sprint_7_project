package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body.String())
	assert.Len(t, strings.Split(responseRecorder.Body.String(), ","), totalCount)
}

func TestMainHandlerWhenCityNotSupported(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=unknown", nil) // Неподдерживаемый город
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)         // Код ответа 400
	assert.Equal(t, "wrong city value\n", responseRecorder.Body.String()) // Сообщение об ошибке
}

func TestMainHandlerWhenCountMissing(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow", nil) // Пропущенный параметр count
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)      // Код ответа 400
	assert.Equal(t, "count missing\n", responseRecorder.Body.String()) // Сообщение об ошибке
}

func TestMainHandlerWhenWrongCountValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=abc&city=moscow", nil) // Неправильное значение count
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)          // Код ответа 400
	assert.Equal(t, "wrong count value\n", responseRecorder.Body.String()) // Сообщение об ошибке
}

func main() {
	http.HandleFunc("/cafe", mainHandle)
	http.ListenAndServe(":8080", nil)
}
