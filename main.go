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


	println("Server läuft auf http://localhost:8080/entrys ...")

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
