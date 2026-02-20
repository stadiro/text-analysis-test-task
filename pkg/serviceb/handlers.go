package serviceb

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// запрос от сервиса A
type AnalyzeRequest struct {
	RequestID string `json:"request_id"`
	Text      string `json:"text"`
}

// ответ для сервиса A
type AnalyzeResponse struct {
	RequestID      string  `json:"request_id"`
	WordCount      int     `json:"word_count"`
	CharCount      int     `json:"char_count"`
	SentenceCount  int     `json:"sentence_count"`
	AverageWordLen float64 `json:"average_word_len"`
}

//http-хендлеры сервиса
type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

//принимает текст и возвращает статистику
func (h *Handlers) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("decode request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	stats := AnalyzeText(req.Text)

	resp := AnalyzeResponse{
		RequestID:      req.RequestID,
		WordCount:      stats.WordCount,
		CharCount:      stats.CharCount,
		SentenceCount:  stats.SentenceCount,
		AverageWordLen: stats.AverageWordLen,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

//возвращает статус сервиса
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

