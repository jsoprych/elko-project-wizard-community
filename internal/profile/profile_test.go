package profile

import (
	"testing"

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

func TestCreateAndGet(t *testing.T) {
	s := newTestStore(t)
	p := &Profile{
		Name: "My Profile", Description: "Test profile",
		ProjectName: "my-project", DirectiveIDs: []string{"d1", "d2"},
		Variables: map[string]interface{}{"max_func_lines": 80},
	}
	created, err := s.Create(p)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Error("expected non-empty ID")
	}
	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != p.Name {
		t.Errorf("Name: got %q want %q", got.Name, p.Name)
	}
	if len(got.DirectiveIDs) != 2 {
		t.Errorf("DirectiveIDs: got %d want 2", len(got.DirectiveIDs))
	}
}

func TestUpdateAndDelete(t *testing.T) {
	s := newTestStore(t)
	created, _ := s.Create(&Profile{Name: "Old Name", DirectiveIDs: []string{"d1"}})

	created.Name = "New Name"
	created.DirectiveIDs = []string{"d2", "d3"}
	if err := s.Update(created); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, _ := s.Get(created.ID)
	if got.Name != "New Name" {
		t.Errorf("Name after update: got %q", got.Name)
	}
	if len(got.DirectiveIDs) != 2 {
		t.Errorf("DirectiveIDs after update: got %d want 2", len(got.DirectiveIDs))
	}
	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	list, _ := s.List()
	if len(list) != 0 {
		t.Errorf("expected empty list after delete, got %d", len(list))
	}
}
