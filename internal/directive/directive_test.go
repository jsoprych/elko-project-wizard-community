package directive

import (
	"fmt"
	"testing"
	"time"

	"elko-project-wizard/internal/db"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	database, err := db.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return NewStore(database.Conn())
}

func TestUpsertAndGet(t *testing.T) {
	s := newTestStore(t)
	d := &Directive{
		ID: "test-directive", Name: "Test", Category: "source-rules",
		Description: "A test directive", Content: "## Test\nDo stuff.",
		Tags: []string{"test"}, Builtin: false, CreatedAt: time.Now(),
	}
	if err := s.Upsert(d); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	got, err := s.Get("test-directive")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != d.Name {
		t.Errorf("Name: got %q want %q", got.Name, d.Name)
	}
}

func TestList(t *testing.T) {
	s := newTestStore(t)
	for i, cat := range []string{"source-rules", "visibility", "source-rules"} {
		s.Upsert(&Directive{
			ID: fmt.Sprintf("d-%s-%d", cat, i),
			Name: cat + " directive", Category: cat,
			Content: "content", CreatedAt: time.Now(),
		})
	}
	all, err := s.List("")
	if err != nil {
		t.Fatalf("List all: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 directives, got %d", len(all))
	}
	filtered, err := s.List("visibility")
	if err != nil {
		t.Fatalf("List filtered: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("expected 1 visibility directive, got %d", len(filtered))
	}
}

func TestDeleteBuiltinBlocked(t *testing.T) {
	s := newTestStore(t)
	s.Upsert(&Directive{
		ID: "builtin-d", Name: "Builtin", Category: "tech-stack",
		Content: "content", Builtin: true, CreatedAt: time.Now(),
	})
	if err := s.Delete("builtin-d"); err == nil {
		t.Error("expected error deleting builtin directive, got nil")
	}
}
