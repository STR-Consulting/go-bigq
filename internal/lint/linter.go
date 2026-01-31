// Package lint provides the core SQL linting logic.
package lint

import (
	"fmt"
	"os"
	"strings"

	"github.com/pacer/go-bigq/bigq"
)

// Result represents a single lint finding.
type Result struct {
	File    string `json:"file"`
	Line    int    `json:"line"`    // 1-based
	Column  int    `json:"column"`  // 1-based
	Level   string `json:"level"`   // "error" or "warning"
	Message string `json:"message"`
}

func (r Result) String() string {
	if r.File != "" && r.Line > 0 {
		return fmt.Sprintf("%s:%d:%d: %s: %s", r.File, r.Line, r.Column, r.Level, r.Message)
	}
	if r.File != "" {
		return fmt.Sprintf("%s: %s: %s", r.File, r.Level, r.Message)
	}
	return fmt.Sprintf("%s: %s", r.Level, r.Message)
}

// Linter validates SQL statements against a catalog.
type Linter struct {
	catalog *bigq.Catalog
}

// New creates a new Linter with the given catalog.
func New(catalog *bigq.Catalog) *Linter {
	return &Linter{catalog: catalog}
}

// LintSQL checks a SQL string (potentially multi-statement) for errors.
func (l *Linter) LintSQL(sql string) []Result {
	statements := splitStatements(sql)
	var results []Result

	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt.text)
		if trimmed == "" || trimmed == ";" {
			continue
		}

		// Skip DECLARE statements — valid BigQuery scripting syntax
		// that ZetaSQL's parser doesn't support.
		upper := strings.ToUpper(trimmed)
		if strings.HasPrefix(upper, "DECLARE ") || strings.HasPrefix(upper, "DECLARE\n") || strings.HasPrefix(upper, "DECLARE\t") {
			continue
		}

		var err error
		if l.catalog != nil {
			err = bigq.AnalyzeStatement(trimmed, l.catalog)
		} else {
			err = bigq.ParseStatement(trimmed)
		}

		if err != nil {
			results = append(results, Result{
				Line:    stmt.startLine,
				Column:  1,
				Level:   "error",
				Message: err.Error(),
			})
		}
	}

	return results
}

// LintFile reads and lints a SQL file.
func (l *Linter) LintFile(path string) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	results := l.LintSQL(string(data))
	for i := range results {
		results[i].File = path
	}
	return results, nil
}

type stmtSpan struct {
	text      string
	startLine int
}

// splitStatements splits SQL on semicolons, tracking line numbers.
func splitStatements(sql string) []stmtSpan {
	var spans []stmtSpan
	line := 1
	start := 0
	startLine := 1
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sql); i++ {
		c := sql[i]

		if c == '\n' {
			line++
			if inLineComment {
				inLineComment = false
			}
			continue
		}

		if inLineComment {
			continue
		}

		if inBlockComment {
			if c == '*' && i+1 < len(sql) && sql[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inSingleQuote {
			if c == '\'' {
				inSingleQuote = false
			} else if c == '\\' {
				i++ // skip escaped char
			}
			continue
		}

		if inDoubleQuote {
			if c == '"' {
				inDoubleQuote = false
			} else if c == '\\' {
				i++
			}
			continue
		}

		if inBacktick {
			if c == '`' {
				inBacktick = false
			}
			continue
		}

		switch c {
		case '\'':
			inSingleQuote = true
		case '"':
			inDoubleQuote = true
		case '`':
			inBacktick = true
		case '-':
			if i+1 < len(sql) && sql[i+1] == '-' {
				inLineComment = true
				i++
			}
		case '/':
			if i+1 < len(sql) && sql[i+1] == '*' {
				inBlockComment = true
				i++
			}
		case ';':
			spans = append(spans, stmtSpan{
				text:      sql[start:i],
				startLine: startLine,
			})
			start = i + 1
			startLine = line
		}
	}

	// Remaining text after last semicolon
	if start < len(sql) {
		remaining := strings.TrimSpace(sql[start:])
		if remaining != "" {
			spans = append(spans, stmtSpan{
				text:      sql[start:],
				startLine: startLine,
			})
		}
	}

	return spans
}
