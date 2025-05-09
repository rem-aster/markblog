package webapp

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"strings"

	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	"encore.app/api"
	"encore.app/api/db"
)

var secrets struct {
	SessionSecret string
}

var store = sessions.NewCookieStore([]byte(secrets.SessionSecret))

func InitService() {
	if secrets.SessionSecret == "" {
		println("Warning: SessionSecret is empty")
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

//go:embed frontend/dist/*
var frontend embed.FS

//encore:api public raw path=/!fallback tag:noclient
func Serve(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	subFS, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestPath := path.Clean(r.URL.Path)
	filePath := strings.TrimPrefix(requestPath, "/")
	if filePath == "" {
		filePath = "index.html"
	}
	file, err := subFS.Open(filePath)
	if err != nil {
		r.URL.Path = "/index.html"
	} else {
		file.Close()
	}
	if strings.HasPrefix(requestPath, "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	}
	http.FileServer(http.FS(subFS)).ServeHTTP(w, r)
}

//encore:api public raw path=/app/auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	username := req.Username
	password := req.Password

	if username == "" || password == "" {
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-]*$`).MatchString(username) {
		http.Error(w, `{"error":"Username has illegal characters"}`, http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, `{"error":"Password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}

	if len(username) < 3 || len(username) > 30 {
		http.Error(w, `{"error":"Username length must be between 3 and 30"}`, http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, `{"error":"Password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}

	exists, err := api.CheckUserExists(r.Context(), username)
	
	if err != nil && !errors.Is(err, sqldb.ErrNoRows) {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	if exists.Exists {
		http.Error(w, `{"error":"Username already taken"}`, http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"Failed to secure password"}`, http.StatusInternalServerError)
		return
	}

	user, err := api.CreateUser(r.Context(), db.CreateUserParams{
		Username:     username,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "markblog")
	if err != nil {
		println("Session get error:", err.Error())
		http.Error(w, `{"error":"Session error"}`, http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID.String()
	session.Values["username"] = user.Username

	if err := session.Save(r, w); err != nil {
		println("Session save error:", err.Error())
		http.Error(w, `{"error":"Failed to save session"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

//encore:api public raw path=/app/auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	username := req.Username
	password := req.Password

	if username == "" || password == "" {
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}
	
	exists, err := api.CheckUserExists(r.Context(), username)
	
	if err != nil && !errors.Is(err, sqldb.ErrNoRows) {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	if !exists.Exists {
		http.Error(w, `{"error":"User does not exist"}`, http.StatusNotFound)
		return
	}

	user, err := api.GetUserByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, sqldb.ErrNoRows) {
			http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, "markblog")
	if err != nil {
		println("Session get error:", err.Error())
		http.Error(w, `{"error":"Session error"}`, http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID.String()
	session.Values["username"] = user.Username

	if err := session.Save(r, w); err != nil {
		println("Session save error:", err.Error())
		http.Error(w, `{"error":"Failed to save session"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

//encore:api public raw path=/app/auth/check
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	session, err := store.Get(r, "markblog")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"error":         "Session error",
		})
		return
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	userID := session.Values["user_id"].(string)
	username := session.Values["username"].(string)
	
	if !ok || !authenticated {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}
	
	exists, err := api.CheckUserExists(r.Context(), username)
	
	if err != nil && !errors.Is(err, sqldb.ErrNoRows) || !exists.Exists{
		session.Values["authenticated"] = false
		session.Values["user_id"] = ""
		session.Values["username"] = ""
		session.Options.MaxAge = -1 // Immediately expire the session

		if err := session.Save(r, w); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Failed to save session",
			})
			return
		}
		http.Error(w, `{"error":"User outdated"}`, http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":       userID,
			"username": username,
		},
	})
}

//encore:api public raw path=/app/auth/logout
func Logout(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	session, err := store.Get(r, "markblog")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session error",
		})
		return
	}

	session.Values["authenticated"] = false
	session.Values["user_id"] = ""
	session.Values["username"] = ""
	session.Options.MaxAge = -1 // Immediately expire the session

	if err := session.Save(r, w); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to save session",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

//encore:api public raw path=/app/post
func Post(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	content := req.Content

	if content == "" || len(content) > 300{
		http.Error(w, `{"error":"Content length is invalid}`, http.StatusBadRequest)
		return
	}
	
	session, err := store.Get(r, "markblog")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session error",
		})
		return
	}
	
	post, err := api.CreatePost(r.Context(), db.CreatePostParams{
		UserID: uuid.Must(uuid.FromString(session.Values["user_id"].(string))),
		Content: content,
	})
	
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": post.ID,
	})
}

//encore:api public raw path=/app/feed
func Feed(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		Offset int32 `json:"offset"`
		Limit int32 `json:"limit"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	offset := req.Offset
	limit := req.Limit
	
	posts, err := api.GetLatestPosts(r.Context(), db.GetLatestPostsParams{
		Offset: offset,
		Limit: limit,
	})
	
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts.Posts,
	})
}

//encore:api public raw path=/app/discussion
func Discussion(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		ID uuid.UUID `json:"id"`
		Limit int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	id := req.ID
	limit := req.Limit
	offset := req.Offset
	
	comments, err := api.GetLatestCommentsForPost(r.Context(), db.GetLatestCommentsForPostParams{
		PostID: id,
		Limit: limit,
		Offset: offset,
	})
	
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"comments": comments.Comments,
	})
}

//encore:api public raw path=/app/comment
func Comment(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		PostID uuid.UUID `json:"post_id"`
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	postID := req.PostID
	content:= req.Content
	
	if content == "" || len(content) > 128{
		http.Error(w, `{"error":"Content length is invalid}`, http.StatusBadRequest)
		return
	}
	
	session, err := store.Get(r, "markblog")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session error",
		})
		return
	}
	
	userID := uuid.Must(uuid.FromString(session.Values["user_id"].(string)))
	
	comment, err := api.CreateComment(r.Context(), db.CreateCommentParams{
		PostID: postID,
		UserID: &userID,
		Content: content,
	})
	
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": comment.ID,
	})
}

//encore:api public raw path=/app/activity
func Activity(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://127.0.0.1:4000" || origin == "http://localhost:4000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		ID uuid.UUID `json:"id"`
		Limit int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	id := req.ID
	limit := req.Limit
	offset := req.Offset
	
	res, err := api.GetLatestUserActivity(r.Context(), db.GetLatestUserActivityParams{
		UserID: id,
		Limit: limit,
		Offset: offset,
	})
	
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"comments": res.Activity,
	})
}