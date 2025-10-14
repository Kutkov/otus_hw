package main

import (
	"log"
	"net/http"
	"os"

	pkgdb "dialog-service/internal/db"
	pkghttp "dialog-service/internal/http"
)

func main() {
	databasePath := os.Getenv("SQLITE_PATH")
	if databasePath == "" {
		databasePath = "./dialog-service.db"
	}
	database, err := pkgdb.Open(databasePath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := pkgdb.Migrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	server := pkghttp.NewServer(database)
	addr := ":8081" // Different port from monolith
	log.Printf("dialog service listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
