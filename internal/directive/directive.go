package directive

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Directive is a single modular AI agent policy block.
type Directive struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Category    string            `json:"category"`
	Description string            `json:"description"`
	Content     string            `json:"content"` // markdown, composed into AGENTS.md
	Tags        []string          `json:"tags"`
	Variables   map[string]VarDef `json:"variables"`
	Builtin     bool              `json:"builtin"`
	CreatedAt   time.Time         `json:"created_at"`
}

// VarDef describes a template variable within a directive.
type VarDef struct {
	Type        string      `json:"type"` // string | int | bool
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
}

// Store manages directives in SQLite.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Upsert inserts or replaces a directive.
func (s *Store) Upsert(d *Directive) error {
	tags, _ := json.Marshal(d.Tags)
	vars, _ := json.Marshal(d.Variables)
	_, err := s.db.Exec(`
		INSERT INTO directives (id, name, category, description, content, tags, variables, builtin, created_at)
		VALUES (?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, category=excluded.category,
			description=excluded.description, content=excluded.content,
			tags=excluded.tags, variables=excluded.variables`,
		d.ID, d.Name, d.Category, d.Description, d.Content,
		string(tags), string(vars), boolToInt(d.Builtin), d.CreatedAt)
	return err
}

// List returns all directives, optionally filtered by category.
func (s *Store) List(category string) ([]*Directive, error) {
	q := `SELECT id,name,category,description,content,tags,variables,builtin,created_at FROM directives`
	args := []interface{}{}
	if category != "" {
		q += ` WHERE category=?`
		args = append(args, category)
	}
	q += ` ORDER BY category, name`

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Directive
	for rows.Next() {
		d, err := scanDirective(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// Get returns a single directive by ID.
func (s *Store) Get(id string) (*Directive, error) {
	row := s.db.QueryRow(`
		SELECT id,name,category,description,content,tags,variables,builtin,created_at
		FROM directives WHERE id=?`, id)
	return scanDirective(row)
}

// Delete removes a non-builtin directive.
func (s *Store) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM directives WHERE id=? AND builtin=0`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("directive %q not found or is builtin", id)
	}
	return nil
}

// Categories returns distinct category names.
func (s *Store) Categories() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT category FROM directives ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanDirective(s scanner) (*Directive, error) {
	var d Directive
	var tagsJSON, varsJSON string
	var builtin int
	err := s.Scan(&d.ID, &d.Name, &d.Category, &d.Description,
		&d.Content, &tagsJSON, &varsJSON, &builtin, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(tagsJSON), &d.Tags)
	json.Unmarshal([]byte(varsJSON), &d.Variables)
	d.Builtin = builtin == 1
	return &d, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
