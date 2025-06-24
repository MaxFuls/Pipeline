package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"pipeline/internal/app"
)

type FormData struct {
	MainMessage      string   `json:"mainMessage"`
	AdditionalFields []string `json:"additionalFields"`
	Timestamp        string   `json:"timestamp"`
}

type Response struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    FormData `json:"data,omitempty"`
}

type HTTPApi struct {
	pipelineService *app.PipelineService
}

func New(pipelineService *app.PipelineService) *HTTPApi {
	return &HTTPApi{
		pipelineService: pipelineService,
	}
}

func (h *HTTPApi) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
		return
	}
}

func (h *HTTPApi) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовки CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Обработка preflight запросов
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Парсим JSON данные
	var formData FormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Ошибка парсинга данных: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Добавляем временную метку
	formData.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	// Логируем полученные данные
	fmt.Printf("📨 Получены данные:\n")
	fmt.Printf("   Основное сообщение: %s\n", formData.MainMessage)
	fmt.Printf("   Дополнительные поля: %v\n", formData.AdditionalFields)
	fmt.Printf("   Время: %s\n", formData.Timestamp)
	fmt.Println("---")

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Получаем список обработчиков из registry
	// Используем дополнительные поля как имена обработчиков
	handlerNames := formData.AdditionalFields
	if len(handlerNames) == 0 {
		// Если дополнительных полей нет, используем дефолтные обработчики
		handlerNames = []string{"validator", "processor", "formatter"}
	}

	fmt.Printf("🔍 Получение обработчиков: %v\n", handlerNames)
	handlers, err := h.pipelineService.GetHandlers(ctx, handlerNames)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Ошибка получения обработчиков: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(handlers) == 0 {
		response := Response{
			Success: false,
			Message: "Не найдено доступных обработчиков",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Обрабатываем главное сообщение через цепочку обработчиков
	fmt.Printf("⚙️ Обработка сообщения через %d обработчиков\n", len(handlers))
	processedMessage, err := h.pipelineService.ProcessMessage(ctx, formData.MainMessage, handlers)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Ошибка обработки сообщения: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Обновляем главное сообщение результатом обработки
	formData.MainMessage = processedMessage

	fmt.Printf("✅ Обработка завершена. Результат: %s\n", processedMessage)

	// Отправляем ответ
	response := Response{
		Success: true,
		Message: "Данные успешно обработаны через цепочку обработчиков!",
		Data:    formData,
	}

	json.NewEncoder(w).Encode(response)
}
