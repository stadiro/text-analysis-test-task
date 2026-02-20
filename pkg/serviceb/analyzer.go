package serviceb

import (
	"regexp"
	"strings"
)

//результаты анализа текста
type TextStats struct {
	WordCount      int     `json:"word_count"`
	CharCount      int     `json:"char_count"`
	SentenceCount  int     `json:"sentence_count"`
	AverageWordLen float64 `json:"average_word_len"`
}

// подсчет слов, символов, предложений и средней длины слова
func AnalyzeText(text string) TextStats {
	stats := TextStats{}
	stats.CharCount = len(text)
	sentenceEnd := regexp.MustCompile(`[.!?]+`)
	sentences := sentenceEnd.Split(text, -1)
	count := 0
	for _, s := range sentences {
		if strings.TrimSpace(s) != "" {
			count++
		}
	}
	stats.SentenceCount = count
	if stats.SentenceCount == 0 && text != "" {
		stats.SentenceCount = 1
	}

	words := strings.Fields(text)
	stats.WordCount = len(words)

	if stats.WordCount > 0 {
		totalLen := 0
		for _, w := range words {
			totalLen += len(w)
		}
		stats.AverageWordLen = float64(totalLen) / float64(stats.WordCount)
	}

	return stats
}

