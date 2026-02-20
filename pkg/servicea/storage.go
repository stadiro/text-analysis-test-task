package servicea

import (
	"sync"
	"time"
)

// статус обработки запроса
type Status string

const (
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

//результаты анализа от сервиса B
type AnalysisResult struct {
	WordCount        int     `json:"word_count"`
	CharCount        int     `json:"char_count"`
	SentenceCount    int     `json:"sentence_count"`
	AverageWordLen   float64 `json:"average_word_len"`
}

//статус запроса и результат
type RequestStatus struct {
	ID        string          `json:"id"`
	Status    Status          `json:"status"`
	Result    *AnalysisResult `json:"result,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

//in-memory хранилище статусов запросов
type Storage struct {
	mu      sync.RWMutex
	records map[string]*RequestStatus
}

// создание нового хранилища
func NewStorage() *Storage {
	return &Storage{
		records: make(map[string]*RequestStatus),
	}
}

// создание новой записи со статусом processing
func (s *Storage) Save(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.records[id] = &RequestStatus{
		ID:        id,
		Status:    StatusProcessing,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// обновление записи результатом анализа
func (s *Storage) UpdateResult(id string, result *AnalysisResult) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.records[id]
	if !ok {
		return false
	}

	record.Status = StatusCompleted
	record.Result = result
	record.UpdatedAt = time.Now()
	return true
}

// помечает запись как failed
func (s *Storage) UpdateFailed(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if record, ok := s.records[id]; ok {
		record.Status = StatusFailed
		record.UpdatedAt = time.Now()
	}
}

func (s *Storage) Get(id string) (*RequestStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.records[id]
	if !ok {
		return nil, false
	}
	return record, true
}

