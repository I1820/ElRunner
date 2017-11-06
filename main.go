package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("GoRunner by Parham Alvani")

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/about", about).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Printf("Listen on %s\n", server.Addr)
	server.ListenAndServe()
}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "18.20 is leaving us")
}
