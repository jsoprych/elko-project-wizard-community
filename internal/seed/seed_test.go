package seed

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"elko-project-wizard/internal/db"
	"elko-project-wizard/internal/directive"
)

func TestLoadBuiltins(t *testing.T) {
	// Create a temp directives directory with two test JSON files
	dir := t.TempDir()
	directives := []directive.Directive{
		{ID: "test-d1", Name: "D1", Category: "source-rules", Content: "## D1"},
		{ID: "test-d2", Name: "D2", Category: "visibility",   Content: "## D2"},
	}
	for _, d := range directives {
		raw, _ := json.Marshal(d)
		sub := filepath.Join(dir, d.Category)
		os.MkdirAll(sub, 0755)
		os.WriteFile(filepath.Join(sub, d.ID+".json"), raw, 0644)
	}

	database, err := db.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()

	if err := LoadBuiltins(database.Conn(), dir); err != nil {
		t.Fatalf("LoadBuiltins: %v", err)
	}

	store := directive.NewStore(database.Conn())
	all, err := store.List("")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 directives, got %d", len(all))
	}

	// Idempotent — running again should not error or duplicate
	if err := LoadBuiltins(database.Conn(), dir); err != nil {
		t.Fatalf("LoadBuiltins (2nd run): %v", err)
	}
	all2, _ := store.List("")
	if len(all2) != 2 {
		t.Errorf("expected still 2 after second load, got %d", len(all2))
	}
}
