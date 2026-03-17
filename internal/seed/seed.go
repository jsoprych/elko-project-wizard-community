// Package seed loads built-in directives from JSON files into the database.
package seed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"elko-project-wizard/internal/directive"
)

// LoadBuiltins walks the directives directory and upserts all .json files.
func LoadBuiltins(db *sql.DB, directivesDir string) error {
	store := directive.NewStore(db)
	count := 0

	err := filepath.Walk(directivesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		var d directive.Directive
		if err := json.Unmarshal(raw, &d); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		d.Builtin = true
		if d.CreatedAt.IsZero() {
			d.CreatedAt = time.Now()
		}
		if err := store.Upsert(&d); err != nil {
			return fmt.Errorf("upsert %s: %w", d.ID, err)
		}
		count++
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("elko-project-wizard: loaded %d built-in directives\n", count)
	return nil
}
