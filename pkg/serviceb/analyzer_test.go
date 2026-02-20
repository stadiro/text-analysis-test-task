package serviceb

import (
	"testing"
)

func TestAnalyzeText(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		wantWords      int
		wantChars      int
		wantSentences  int
		wantAvgWordLen float64
	}{
		{
			name:           "simple",
			text:           "Hello world",
			wantWords:      2,
			wantChars:      11,
			wantSentences:  1,
			wantAvgWordLen: 5.0,
		},
		{
			name:           "multiple sentences",
			text:           "First. Second! Third?",
			wantWords:      3,
			wantChars:      21,
			wantSentences:  3,
			wantAvgWordLen: 19.0 / 3.0, // First.=6, Second!=7, Third?=6
		},
		{
			name:           "empty",
			text:           "",
			wantWords:      0,
			wantChars:      0,
			wantSentences:  0,
			wantAvgWordLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnalyzeText(tt.text)
			if got.WordCount != tt.wantWords {
				t.Errorf("WordCount = %d, want %d", got.WordCount, tt.wantWords)
			}
			if got.CharCount != tt.wantChars {
				t.Errorf("CharCount = %d, want %d", got.CharCount, tt.wantChars)
			}
			if got.SentenceCount != tt.wantSentences {
				t.Errorf("SentenceCount = %d, want %d", got.SentenceCount, tt.wantSentences)
			}
			if got.AverageWordLen != tt.wantAvgWordLen {
				t.Errorf("AverageWordLen = %f, want %f", got.AverageWordLen, tt.wantAvgWordLen)
			}
		})
	}
}

