package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Entry struct {
	Title string `json:"title"`
	Code  string `json:"code"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	entries := []Entry{
		{"Test Ausgeben", "fmt.Print(`TEXT`)"},
	}

	mux.HandleFunc("GET /entrys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})

	mux.HandleFunc("POST /entrys", func(w http.ResponseWriter, r *http.Request) {
		var entry Entry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		entries = append(entries, entry)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entry)

	})

	mux.HandleFunc("DELETE /entrys/{title}", func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")

		for i, entry := range entries {
			if entry.Title == title {
				entries = append(entries[:i], entries[i+1:]...)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
			}
		}
	})

	println("Server läuft auf http://localhost:8080/entrys ...")

	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}
