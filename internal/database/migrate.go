package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyMigrations(ctx context.Context, db *sql.DB, schemaDir string) error {
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("read schema dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		files = append(files, filepath.Join(schemaDir, name))
	}
	sort.Strings(files) // 001_..., 002_..., 003_..., 004_...

	for _, f := range files {
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		// naive: run the Up section (before -- +goose Down)
		parts := strings.Split(string(sqlBytes), "-- +goose Down")
		up := parts[0]
		if _, err := db.ExecContext(ctx, up); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}
	}
	return nil
}
