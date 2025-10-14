package main

import (
	"log"
	"net/http"
	"os"

	pkgdb "monolith/internal/db"
	pkghttp "monolith/internal/interfaces/http"
)

func main() {
	databasePath := os.Getenv("SQLITE_PATH")
	if databasePath == "" {
		databasePath = "./monolith.db"
	}
	database, err := pkgdb.Open(databasePath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := pkgdb.Migrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	server := pkghttp.NewServer(database)
	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
