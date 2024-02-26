package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kiquetal/boot-dev/build/server/internal"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
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
	_, err := w.Write([]byte("Hits reset"))
	if err != nil {
		return
	}
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

func (cfg *apiConfig) validateChirpFunc(r *http.Request) (bool, string, error) {
	fmt.Printf("Validate chirp: %v\n", r)
	type bodyJSON struct {
		Body string `json:"body"`
	}
	fmt.Printf("ValidateChripRequest: %v\n", r)
	decoder := json.NewDecoder(r.Body)
	var body bodyJSON
	err2 := decoder.Decode(&body)
	if body.Body == "" {
		return false, "", fmt.Errorf("Invalid request payload")
	}
	if err2 != nil {
		return false, "", fmt.Errorf("Invalid request payload")
	}
	if len(body.Body) > 140 {
		return false, "", fmt.Errorf("Chirp is too long")

	}
	alteredBody := processAndReplaceBadWords(body.Body)
	if alteredBody != body.Body && len(alteredBody) > 0 {
		fmt.Printf("Altered body: %s\n", alteredBody)

		return true, strings.TrimSpace(alteredBody), nil
	} else {

		return true, strings.TrimSpace(body.Body), nil

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
	alteredBody := processAndReplaceBadWords(body.Body)
	if alteredBody != body.Body && len(alteredBody) > 0 {
		fmt.Printf("Altered body: %s\n", alteredBody)
		respondWithJSON(w, http.StatusOK, struct {
			CleanBody string `json:"cleaned_body"`
		}{
			CleanBody: strings.TrimSpace(alteredBody),
		})
	} else {
		respondWithJSON(w, http.StatusOK, struct {
			CleanBody string `json:"cleaned_body"`
		}{
			CleanBody: strings.TrimSpace(body.Body),
		})
	}

}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Create chirp: %v\n", r)
	_, body, err := cfg.validateChirpFunc(r)
	fmt.Printf("RequestCreateChirps: %v\n", r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("Invalid request payload").Error())
		return
	}
	chirp, err := cfg.DB.CreateChirp(body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, struct {
		Body string `json:"body"`
		Id   int    `json:"id"`
	}{
		Body: chirp.Body,
		Id:   chirp.Id,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err2 := encoder.Encode(payload)
	if err2 != nil {
		return
	}
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	//sort by id
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id

	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Printf("ID: %v\n", id)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.DB.GetChirp(idInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userJSON struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	var user userJSON
	err2 := decoder.Decode(&user)
	if user.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return

	}
	if err2 != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return

	}
	userCreated, err := cfg.DB.CreateUser(user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}
	respondWithJSON(w, http.StatusCreated, struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}{
		Email: userCreated.Email,
		Id:    userCreated.Id,
	})
}

func processAndReplaceBadWords(body string) string {

	var badWords = []string{"kerfuffle", "sharbert", "fornax"}
	var newBody string
	var containBadWord bool

	for _, word := range strings.Split(body, " ") {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				containBadWord = true
				break
			}

		}
		if containBadWord {
			newBody += "**** "
			containBadWord = false
		} else {
			newBody += word + " "
		}
	}

	return newBody
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
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

	db, err := database.NewDB("database.json")
	if err != nil {
		panic(err)
	}
	api.DB = db

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
	//	rapi.Post("/validate_chirp", api.validateChirp)
	rapi.Post("/chirps", api.createChirp)
	rapi.Get("/chirps", api.getAllChirps)
	rapi.Get("/chirps/{id}", api.getChirpByID)
	rapi.Post("/users", api.createUser)
	r.Mount("/api", rapi)

	corsR := addCORSHeaders(r)

	// Create a new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsR,
	}

	// Start the server

	println("Server is listening on", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	// Print listener address

}
