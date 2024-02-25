package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
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

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	//write OK
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do not increment if is options
		if r.Method != "OPTIONS" {
			cfg.fileserverHits++
		}
		fmt.Printf("Fileserver hits: %d\n", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) newHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileserverHits)))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset"))
}

type Data struct {
	FileserverHits string
}

func (cfg *apiConfig) templateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// ignore options and is not metric endpoint

	fmt.Printf("Fileserver hits: %d\n", cfg.fileserverHits)

	tmpl, err := template.New("admin").Parse(`
			<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited {{.FileserverHits}} times!</p>
			</body>
			</html>
		`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	d := Data{FileserverHits: strconv.Itoa(cfg.fileserverHits)}
	err = tmpl.Execute(w, d)
	if err != nil {
		return
	}

}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Request: %v\n", r)
	type bodyJSON struct {
		Body string `json:"body"`
	}

	type successJSON struct {
		Valid bool `json:"valid"`
	}
	fmt.Printf("Request: %v\n", r)
	decoder := json.NewDecoder(r.Body)
	var body bodyJSON
	err2 := decoder.Decode(&body)
	if body.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return

	}
	if err2 != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return

	}

	respondWithJSON(w, http.StatusOK, successJSON{Valid: true})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err2 := encoder.Encode(payload)
	if err2 != nil {
		return
	}
}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err2 := encoder.Encode(struct {
		Error string `json:"err"`
	}{Error: msg})
	if err2 != nil {
		return
	}
}

var api = &apiConfig{}

func main() {

	r := chi.NewRouter()
	fsHandler := api.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	rapi := chi.NewRouter()
	radmin := chi.NewRouter()
	// Wrap the mux with the CORS middleware
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)
	radmin.Get("/metrics", api.templateAdmin)
	r.Mount("/admin", radmin)
	rapi.Get("/healthz", myHandler)
	rapi.Get("/metrics", api.newHandler)
	rapi.Get("/reset", api.reset)
	rapi.Post("/validate_chirp", api.validateChirp)
	r.Mount("/api", rapi)

	corsR := addCORSHeaders(r)

	// Create a new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsR,
	}

	// Start the server

	println("Server is listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	// Print listener address

}
