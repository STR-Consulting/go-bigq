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

func TestLintSQL_SkipsDeclare(t *testing.T) {
	l := &Linter{} // parse-only, no catalog

	tests := []struct {
		name string
		sql  string
	}{
		{"declare with default", "DECLARE run_date DATE DEFAULT CURRENT_DATE();\nSELECT 1"},
		{"declare no default", "DECLARE inserted_rows INT64;\nSELECT 1"},
		{"declare lowercase", "declare x INT64;\nSELECT 1"},
		{"declare mixed case", "Declare x INT64;\nSELECT 1"},
		{"multiple declares", "DECLARE a INT64;\nDECLARE b STRING;\nSELECT 1"},
		{"declare only", "DECLARE x INT64"},
		{"declare with tab", "DECLARE\tx INT64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := l.LintSQL(tt.sql)
			if len(results) != 0 {
				t.Errorf("LintSQL(%q) returned %d errors, want 0", tt.sql, len(results))
				for _, r := range results {
					t.Logf("  %s", r)
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
