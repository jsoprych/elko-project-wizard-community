package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Profile is a saved composition of directives.
type Profile struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ProjectName  string                 `json:"project_name"`
	DirectiveIDs []string               `json:"directive_ids"`
	Variables    map[string]interface{} `json:"variables"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// Store manages profiles in SQLite.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Create inserts a new profile and returns it with a generated ID.
func (s *Store) Create(p *Profile) (*Profile, error) {
	p.ID = uuid.New().String()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	vars, _ := json.Marshal(p.Variables)
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO profiles (id,name,description,project_name,variables,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?)`,
		p.ID, p.Name, p.Description, p.ProjectName, string(vars), p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := insertDirectives(tx, p.ID, p.DirectiveIDs); err != nil {
		return nil, err
	}
	return p, tx.Commit()
}

// Update replaces a profile's directive list and metadata.
func (s *Store) Update(p *Profile) error {
	p.UpdatedAt = time.Now()
	vars, _ := json.Marshal(p.Variables)
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		UPDATE profiles SET name=?,description=?,project_name=?,variables=?,updated_at=?
		WHERE id=?`,
		p.Name, p.Description, p.ProjectName, string(vars), p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("profile %q not found", p.ID)
	}
	if _, err := tx.Exec(`DELETE FROM profile_directives WHERE profile_id=?`, p.ID); err != nil {
		return err
	}
	if err := insertDirectives(tx, p.ID, p.DirectiveIDs); err != nil {
		return err
	}
	return tx.Commit()
}

// Get returns a profile by ID including its directive IDs.
func (s *Store) Get(id string) (*Profile, error) {
	row := s.db.QueryRow(`
		SELECT id,name,description,project_name,variables,created_at,updated_at
		FROM profiles WHERE id=?`, id)
	p, err := scanProfile(row)
	if err != nil {
		return nil, err
	}
	p.DirectiveIDs, err = s.directiveIDs(id)
	return p, err
}

// List returns all profiles.
func (s *Store) List() ([]*Profile, error) {
	rows, err := s.db.Query(`
		SELECT id,name,description,project_name,variables,created_at,updated_at
		FROM profiles ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Profile
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		p.DirectiveIDs, _ = s.directiveIDs(p.ID)
		out = append(out, p)
	}
	return out, rows.Err()
}

// Delete removes a profile.
func (s *Store) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM profiles WHERE id=?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("profile %q not found", id)
	}
	return nil
}

func (s *Store) directiveIDs(profileID string) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT directive_id FROM profile_directives
		WHERE profile_id=? ORDER BY sort_order`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func insertDirectives(tx *sql.Tx, profileID string, ids []string) error {
	for i, did := range ids {
		_, err := tx.Exec(`
			INSERT INTO profile_directives (profile_id,directive_id,sort_order)
			VALUES (?,?,?)`, profileID, did, i)
		if err != nil {
			return err
		}
	}
	return nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanProfile(s scanner) (*Profile, error) {
	var p Profile
	var varsJSON string
	err := s.Scan(&p.ID, &p.Name, &p.Description, &p.ProjectName,
		&varsJSON, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(varsJSON), &p.Variables)
	if p.Variables == nil {
		p.Variables = map[string]interface{}{}
	}
	return &p, nil
}
