package zetasql_test

import (
	"testing"

	"github.com/pacer/go-bigq/zetasql"
)

func TestParseStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"simple select", "SELECT 1", false},
		{"select with alias", "SELECT 1 + 2 AS result", false},
		{"select star", "SELECT * FROM t", false},
		{"syntax error", "SELECT * FORM t", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := zetasql.ParseStatement(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatement(%q) error = %v, wantErr %v", tt.sql, err, tt.wantErr)
			}
		})
	}
}

func TestAnalyzeStatement(t *testing.T) {
	cat, err := zetasql.NewCatalog("test")
	if err != nil {
		t.Fatalf("NewCatalog: %v", err)
	}
	defer cat.Close()

	err = cat.AddTable("my_table", []zetasql.ColumnDef{
		{Name: "id", TypeName: "INT64"},
		{Name: "name", TypeName: "STRING"},
		{Name: "created_at", TypeName: "TIMESTAMP"},
	})
	if err != nil {
		t.Fatalf("AddTable: %v", err)
	}

	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"valid select", "SELECT id, name FROM my_table", false},
		{"select star", "SELECT * FROM my_table", false},
		{"bad column", "SELECT nonexistent FROM my_table", true},
		{"bad table", "SELECT 1 FROM no_such_table", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := zetasql.AnalyzeStatement(tt.sql, cat)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeStatement(%q) error = %v, wantErr %v", tt.sql, err, tt.wantErr)
			}
		})
	}
}
