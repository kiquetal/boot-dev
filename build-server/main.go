package main

import (
	"net/http"
)

// Middleware function to add CORS headers
func addCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add headers
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Wrap the mux with the CORS middleware
	corsMux := addCORSHeaders(mux)

	// Create a new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	// Start the server
	server.ListenAndServe()
}

