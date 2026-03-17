package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"elko-project-wizard/internal/ai"
	"elko-project-wizard/internal/db"
	"elko-project-wizard/internal/directive"
	"elko-project-wizard/internal/generator"
	"elko-project-wizard/internal/profile"
)

// Server is the HTTP server for elko-project-wizard.
type Server struct {
	db        *db.DB
	dirStore  *directive.Store
	profStore *profile.Store
	gen       *generator.Generator
	mux       *http.ServeMux
}

func New(database *db.DB) *Server {
	s := &Server{
		db:        database,
		dirStore:  directive.NewStore(database.Conn()),
		profStore: profile.NewStore(database.Conn()),
		gen:       generator.New(database.Conn()),
		mux:       http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	log.Printf("elko-project-wizard listening on http://%s", addr)
	return http.ListenAndServe(addr, s)
}

// json helpers
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *Server) routes() {
	// Static web UI
	s.mux.Handle("/", http.FileServer(http.Dir("web")))

	// Directives
	s.mux.HandleFunc("/api/directives", s.handleDirectives)
	s.mux.HandleFunc("/api/directives/", s.handleDirective)

	// Profiles
	s.mux.HandleFunc("/api/profiles", s.handleProfiles)
	s.mux.HandleFunc("/api/profiles/", s.handleProfile)

	// Generate
	s.mux.HandleFunc("/api/generate", s.handleGenerate)

	// AI prompt sandwich
	s.mux.HandleFunc("/api/prompt", s.handlePrompt)

	// Generation history
	s.mux.HandleFunc("/api/generations", s.handleGenerations)
}

// --- Directives ---

func (s *Server) handleDirectives(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cat := r.URL.Query().Get("category")
		dirs, err := s.dirStore.List(cat)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, dirs)

	case http.MethodPost:
		var d directive.Directive
		if err := decodeJSON(r, &d); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		d.Builtin = false
		if err := s.dirStore.Upsert(&d); err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 201, d)

	default:
		writeError(w, 405, "method not allowed")
	}
}

func (s *Server) handleDirective(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/directives/"):]
	switch r.Method {
	case http.MethodGet:
		d, err := s.dirStore.Get(id)
		if err != nil {
			writeError(w, 404, fmt.Sprintf("directive %q not found", id))
			return
		}
		writeJSON(w, 200, d)

	case http.MethodDelete:
		if err := s.dirStore.Delete(id); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		w.WriteHeader(204)

	default:
		writeError(w, 405, "method not allowed")
	}
}

// --- Profiles ---

func (s *Server) handleProfiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		profiles, err := s.profStore.List()
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, profiles)

	case http.MethodPost:
		var p profile.Profile
		if err := decodeJSON(r, &p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		created, err := s.profStore.Create(&p)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 201, created)

	default:
		writeError(w, 405, "method not allowed")
	}
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/profiles/"):]
	switch r.Method {
	case http.MethodGet:
		p, err := s.profStore.Get(id)
		if err != nil {
			writeError(w, 404, err.Error())
			return
		}
		writeJSON(w, 200, p)

	case http.MethodPut:
		var p profile.Profile
		if err := decodeJSON(r, &p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		p.ID = id
		if err := s.profStore.Update(&p); err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, p)

	case http.MethodDelete:
		if err := s.profStore.Delete(id); err != nil {
			writeError(w, 404, err.Error())
			return
		}
		w.WriteHeader(204)

	default:
		writeError(w, 405, "method not allowed")
	}
}

// --- Generate ---

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}
	var req struct {
		ProfileID   string `json:"profile_id"`
		ProjectName string `json:"project_name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	result, err := s.gen.Generate(req.ProfileID, req.ProjectName)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, result.ProjectName))
	w.WriteHeader(200)
	w.Write(result.ZipBytes)
}

// --- AI Prompt ---

func (s *Server) handlePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, 405, "method not allowed")
		return
	}
	var req ai.PromptRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, 400, err.Error())
		return
	}

	// Try API first (private edition); fall back to sandwich (community)
	result, err := ai.GenerateViaAPI(&req)
	if err != nil {
		// Fall back: return the prompt sandwich for manual paste
		sandwich, serr := ai.BuildPromptSandwich(&req)
		if serr != nil {
			writeError(w, 500, serr.Error())
			return
		}
		writeJSON(w, 200, map[string]interface{}{
			"mode":   "sandwich",
			"prompt": sandwich,
			"hint":   "Paste this into Claude, ChatGPT, Grok, or DeepSeek — then save the returned JSON as a custom directive.",
		})
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"mode":   "api",
		"result": result,
	})
}

// --- Generations ---

func (s *Server) handleGenerations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, 405, "method not allowed")
		return
	}
	rows, err := s.db.Conn().Query(`
		SELECT id, profile_id, project_name, files, created_at
		FROM generations ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var gens []map[string]interface{}
	for rows.Next() {
		var id, profID, projName, files string
		var createdAt interface{}
		rows.Scan(&id, &profID, &projName, &files, &createdAt)
		gens = append(gens, map[string]interface{}{
			"id": id, "profile_id": profID,
			"project_name": projName, "files": files, "created_at": createdAt,
		})
	}
	writeJSON(w, 200, gens)
}
