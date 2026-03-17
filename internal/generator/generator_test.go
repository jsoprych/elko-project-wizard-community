package generator

import (
	"archive/zip"
	"bytes"
	"testing"
	"time"

	"elko-project-wizard/internal/db"
	"elko-project-wizard/internal/directive"
	"elko-project-wizard/internal/profile"
)

func setupTest(t *testing.T) (*Generator, *directive.Store, *profile.Store) {
	t.Helper()
	database, err := db.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	conn := database.Conn()
	return New(conn), directive.NewStore(conn), profile.NewStore(conn)
}

func seedDirective(t *testing.T, ds *directive.Store, id, cat, content string) {
	t.Helper()
	ds.Upsert(&directive.Directive{
		ID: id, Name: id, Category: cat,
		Content: content, Builtin: true, CreatedAt: time.Now(),
	})
}

func TestGenerateProducesZip(t *testing.T) {
	g, ds, ps := setupTest(t)

	seedDirective(t, ds, "stack-go", "tech-stack", "## Go Stack\nUse Go.")
	seedDirective(t, ds, "rules-minimal-deps", "source-rules", "## Dependencies\nMinimize deps.")

	p, err := ps.Create(&profile.Profile{
		Name:         "Test Profile",
		ProjectName:  "test-project",
		DirectiveIDs: []string{"stack-go", "rules-minimal-deps"},
	})
	if err != nil {
		t.Fatalf("Create profile: %v", err)
	}

	result, err := g.Generate(p.ID, "")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if len(result.ZipBytes) == 0 {
		t.Fatal("expected non-empty zip")
	}

	// Verify zip contains expected files
	zr, err := zip.NewReader(bytes.NewReader(result.ZipBytes), int64(len(result.ZipBytes)))
	if err != nil {
		t.Fatalf("read zip: %v", err)
	}
	fileMap := map[string]bool{}
	for _, f := range zr.File {
		fileMap[f.Name] = true
	}
	for _, expected := range []string{"test-project/AGENTS.md", "test-project/CLAUDE.md", "test-project/.gitignore"} {
		if !fileMap[expected] {
			t.Errorf("missing file in zip: %s", expected)
		}
	}
}

func TestDockerFilesIncludedWhenDirectivePresent(t *testing.T) {
	g, ds, ps := setupTest(t)

	seedDirective(t, ds, "stack-go", "tech-stack", "## Go")
	seedDirective(t, ds, "docker-go", "docker", "## Docker\nUse Docker.")

	p, _ := ps.Create(&profile.Profile{
		Name: "Docker Profile", ProjectName: "docker-project",
		DirectiveIDs: []string{"stack-go", "docker-go"},
	})

	result, err := g.Generate(p.ID, "")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	zr, _ := zip.NewReader(bytes.NewReader(result.ZipBytes), int64(len(result.ZipBytes)))
	fileMap := map[string]bool{}
	for _, f := range zr.File {
		fileMap[f.Name] = true
	}
	if !fileMap["docker-project/Dockerfile"] {
		t.Error("expected Dockerfile in zip when docker directive present")
	}
	if !fileMap["docker-project/docker-compose.yml"] {
		t.Error("expected docker-compose.yml in zip when docker directive present")
	}
}
