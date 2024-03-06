package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kiquetal/boot-dev/build/server/internal"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Middleware function to add CORS headers
func addCORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add headers
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// Call the next handler

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func myHandler(w http.ResponseWriter, _ *http.Request) {
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
	Secret         string
	POLKA          string
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
func (cfg *apiConfig) middlewareKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check Authorization has ApiKey
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "No token found")
			return
		}
		apiKey = strings.Replace(apiKey, "ApiKey ", "", 1)
		fmt.Printf("Apikey: %s\n", apiKey)
		if apiKey != cfg.POLKA {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) newHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileserverHits)))
	if err != nil {
		return
	}
}

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
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

func (cfg *apiConfig) templateAdmin(w http.ResponseWriter, _ *http.Request) {
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
		return false, "", fmt.Errorf("invalid request payload")
	}
	if err2 != nil {
		return false, "", fmt.Errorf("invalid request payload")
	}
	if len(body.Body) > 140 {
		return false, "", fmt.Errorf("chirp is too long")

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
		respondWithError(w, http.StatusBadRequest, errors.New("invalid request payload").Error())
		return
	}
	subId := r.Context().Value("userId").(string)
	intSubID, err := strconv.Atoi(subId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("SubId: %v\n", subId)

	chirp, err := cfg.DB.CreateChirp(body, intSubID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, struct {
		AuthorId int    `json:"author_id"`
		Body     string `json:"body"`
		Id       int    `json:"id"`
	}{
		Body:     chirp.Body,
		Id:       chirp.Id,
		AuthorId: chirp.AuthorId,
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
	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		authorIdInt, err := strconv.Atoi(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		var chirpsByUser []database.Chirp
		for _, chirp := range chirps {
			if chirp.AuthorId == authorIdInt {
				chirpsByUser = append(chirpsByUser, chirp)
			}
		}
		chirps = chirpsByUser

	}
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
		Email    string `json:"email"`
		Password string `json:"password"`
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
	userCreated, err := cfg.DB.CreateUser(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}
	respondWithJSON(w, http.StatusCreated, struct {
		Email       string `json:"email"`
		Id          int    `json:"id"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		Email:       userCreated.Email,
		Id:          userCreated.Id,
		IsChirpyRed: userCreated.IsChirpyRed,
	})
}
func (cfg *apiConfig) updateRoute(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Update route: %v\n", r)
	type userJSON struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var userInput userJSON
	err2 := decoder.Decode(&userInput)
	if err2 != nil {

		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//read from context
	userId := r.Context().Value("userId").(string)
	issuer := r.Context().Value("issuer").(string)
	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Invalid issuer")
		return

	}
	id, e := strconv.Atoi(userId)
	if e != nil {
		respondWithError(w, http.StatusInternalServerError, e.Error())
		return

	}
	user, e := cfg.DB.GetUser(id)
	fmt.Printf("Obtainig user by header id %v", user)
	if e != nil {
		respondWithError(w, http.StatusInternalServerError, e.Error())
		return

	}

	user, e = cfg.DB.UpdateUser(database.User{
		Id:       user.Id,
		Email:    userInput.Email,
		Password: userInput.Password,
	})
	if e != nil {
		respondWithError(w, http.StatusInternalServerError, e.Error())
		return

	}
	respondWithJSON(w, http.StatusOK, struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}{
		Email: user.Email,
		Id:    user.Id,
	})

}

func (cfg *apiConfig) generateJWT(subjectID string, issuer string) (string, error) {
	fmt.Printf("Subject: %s\n", subjectID)
	fmt.Printf("Issuer: %s\n", issuer)
	const ttlForToken = 24 * 60 * 60
	const ttlForRefreshToken = 24 * 60 * 60 * 60
	var ttl int
	if issuer == "chirpy-access" {
		ttl = ttlForToken
	} else {
		ttl = ttlForRefreshToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subjectID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(ttl))),
	})
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userJSON struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		//	ExpiresInSeconds int    `json:"expires_in_seconds"`
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
	userLogged, err := cfg.DB.Login(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	jwtToken, err := cfg.generateJWT(strconv.Itoa(userLogged.Id), "chirpy-access")

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refresh, err := cfg.generateJWT(strconv.Itoa(userLogged.Id), "chirpy-refresh")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}

	respondWithJSON(w, http.StatusOK, struct {
		Email        string `json:"email"`
		Id           int    `json:"id"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		Email:        userLogged.Email,
		Id:           userLogged.Id,
		Token:        jwtToken,
		RefreshToken: refresh,
		IsChirpyRed:  userLogged.IsChirpyRed,
	})
}

func processAndReplaceBadWords(body string) string {

	var badWords = []string{"kerfuffle", "sharer", "fornax"}
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

func (cfg *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the header
		fmt.Printf("Middleware auth: %v\n", r)
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "No token found")
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		fmt.Printf("Token: %s\n", tokenString)
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(api.Secret), nil

		})

		if err != nil {
			fmt.Printf("Error:11111 %v\n", err)
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return

		}
		if !token.Valid {
			fmt.Printf("Token invalid!!!!: %v\n", token)
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return

		}

		subId := token.Claims.(*jwt.RegisteredClaims).Subject
		issuer := token.Claims.(*jwt.RegisteredClaims).Issuer
		ctx := context.WithValue(r.Context(), "userId", subId)
		ctx = context.WithValue(ctx, "issuer", issuer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {

	//check if issuer is chirpy-refresh

	issuer := r.Context().Value("issuer").(string)
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Invalid issuer")
		return
	}
	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.Replace(refreshToken, "Bearer ", "", 1)
	isRevoked, err := cfg.DB.IsRevokedToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	generateJWT, err := cfg.generateJWT(r.Context().Value("userId").(string), "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: generateJWT,
	})

}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {

	issuer := r.Context().Value("issuer").(string)
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Invalid issuer")
		return
	}
	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.Replace(refreshToken, "Bearer ", "", 1)
	_, err := cfg.DB.SaveRevokedToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})

}

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Delete chirp: %v\n", r)
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	userId := r.Context().Value("userId").(string)
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(idInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	fmt.Printf("Chirp: %+v\n", chirp)
	if chirp.AuthorId != userIdInt {
		respondWithError(w, http.StatusForbidden, "Invalid user")
		return
	}

	fmt.Printf("Id  of the chrip i want to delete, %v", idInt)
	fmt.Printf("Id of the user %v", userIdInt)
	fmt.Printf("Chirp: %+v\n", chirp)
	err = cfg.DB.DeleteChrip(idInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})
}

func (cfg *apiConfig) webhookPolka(writer http.ResponseWriter, request *http.Request) {

	type bodyJsonWebhook struct {
		Data struct {
			User_id int `json:"user_id"`
		}
		Event string `json:"event"`
	}
	decoder := json.NewDecoder(request.Body)
	var body bodyJsonWebhook
	err2 := decoder.Decode(&body)
	if err2 != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if body.Event == "user.upgraded" {

		_, err := cfg.DB.CreatePolkaWebhook(body.Data.User_id)
		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, err.Error())
			return

		}
	}
	respondWithJSON(writer, http.StatusOK, struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})
}

var api = &apiConfig{}

func main() {

	db, err := database.NewDB("database.json")
	err = godotenv.Load()
	if err != nil {
		return
	}
	if err != nil {
		panic(err)
	}
	api.DB = db
	api.Secret = os.Getenv("JWT_SECRET")
	api.POLKA = os.Getenv("POLKA_KEY")
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
	rapi.With(api.middlewareAuth).Post("/chirps", api.createChirp)
	rapi.Get("/chirps", api.getAllChirps)
	rapi.Get("/chirps/{id}", api.getChirpByID)
	rapi.With(api.middlewareAuth).Delete("/chirps/{id}", api.DeleteChirp)
	rapi.Post("/users", api.createUser)
	rapi.Post("/polka/webhooks", api.webhookPolka)
	rapi.With(api.middlewareAuth).Put("/users", api.updateRoute)
	rapi.With(api.middlewareAuth).Post("/refresh", api.refreshToken)
	rapi.With(api.middlewareAuth).Post("/revoke", api.revokeToken)
	rapi.Post("/login", api.login)

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
