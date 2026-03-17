package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"elko-project-wizard/internal/db"
	"elko-project-wizard/internal/seed"
	"elko-project-wizard/internal/server"
)

var version = "0.1.0"

func main() {
	port := flag.String("port", "8080", "HTTP port")
	dataDir := flag.String("data", defaultDataDir(), "Data directory for SQLite DB")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("elko-project-wizard v%s\n", version)
		os.Exit(0)
	}

	dbPath := filepath.Join(*dataDir, "wizard.db")
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	// Load built-in directives from directives/ directory
	if err := seed.LoadBuiltins(database.Conn(), "directives"); err != nil {
		log.Printf("warn: load builtins: %v", err)
	}

	srv := server.New(database)
	log.Fatal(srv.Start("0.0.0.0:" + *port))
}

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".data"
	}
	return filepath.Join(home, ".elko-project-wizard")
}
