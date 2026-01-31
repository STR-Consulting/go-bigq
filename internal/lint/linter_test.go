package lint

import (
	"testing"
)

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		count int
	}{
		{"single", "SELECT 1", 1},
		{"two", "SELECT 1; SELECT 2", 2},
		{"trailing semi", "SELECT 1;", 1},
		{"empty between", "SELECT 1;; SELECT 2", 3},
		{"with comments", "-- comment\nSELECT 1;\n/* block */\nSELECT 2", 2},
		{"string with semi", "SELECT 'a;b'", 1},
		{"backtick with semi", "SELECT `a;b`", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spans := splitStatements(tt.sql)
			if len(spans) != tt.count {
				t.Errorf("splitStatements(%q) = %d statements, want %d", tt.sql, len(spans), tt.count)
				for i, s := range spans {
					t.Logf("  [%d] line=%d text=%q", i, s.startLine, s.text)
				}
			}
		})
	}
}

func TestSplitStatementsLineTracking(t *testing.T) {
	sql := "SELECT 1;\n\nSELECT 2;\nSELECT 3"
	spans := splitStatements(sql)
	if len(spans) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(spans))
	}

	expected := []int{1, 1, 3}
	for i, span := range spans {
		if span.startLine != expected[i] {
			t.Errorf("statement %d: startLine = %d, want %d", i, span.startLine, expected[i])
		}
	}
}
