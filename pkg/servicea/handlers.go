package servicea

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

//запрос от пользователя
type TextRequest struct {
	Text string `json:"text"`
}

//ответ пользователю
type TextResponse struct {
	ID string `json:"id"`
}

//http-хендлеры Сервиса A
type Handlers struct {
	storage *Storage
	client  *Client
}

// создание новых хендлеров
func NewHandlers(storage *Storage, client *Client) *Handlers {
	return &Handlers{
		storage: storage,
		client:  client,
	}
}

//функция принимающая текст для анализа
func (h *Handlers) HandlePostText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("decode request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// валидация
	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	h.storage.Save(id)

	// отправка в сервис B в отдельной горутине
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := h.client.SendForAnalysis(ctx, id, req.Text)
		if err != nil {
			log.Printf("send to service B failed: %v", err)
			h.storage.UpdateFailed(id)
			return
		}

		result := &AnalysisResult{
			WordCount:      resp.WordCount,
			CharCount:      resp.CharCount,
			SentenceCount:  resp.SentenceCount,
			AverageWordLen: resp.AverageWordLen,
		}
		h.storage.UpdateResult(id, result)
		log.Printf("request %s completed", id)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(TextResponse{ID: id})
}

//возвращение статус по id
func (h *Handlers) HandleGetStatus(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	record, ok := h.storage.Get(id)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

//возвращение статуса сервиса
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
