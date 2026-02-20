package servicea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// запрос к сервису B
type AnalyzeRequest struct {
	RequestID string `json:"request_id"`
	Text      string `json:"text"`
}

// ответ от сервиса B
type AnalyzeResponse struct {
	RequestID      string  `json:"request_id"`
	WordCount      int     `json:"word_count"`
	CharCount      int     `json:"char_count"`
	SentenceCount  int     `json:"sentence_count"`
	AverageWordLen float64 `json:"average_word_len"`
}

// http-клиент для сервиса B
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// создание нового клиента для сервиса B
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// отпрвка текста в сервис B на анализ
func (c *Client) SendForAnalysis(ctx context.Context, requestID, text string) (*AnalyzeResponse, error) {
	reqBody := AnalyzeRequest{
		RequestID: requestID,
		Text:      text,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/api/v1/analyze"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service B returned status %d: %s", resp.StatusCode, string(body))
	}

	var analyzeResp AnalyzeResponse
	if err := json.Unmarshal(body, &analyzeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &analyzeResp, nil
}
