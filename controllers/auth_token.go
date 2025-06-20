package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"main/repositories/auth_token"
	"main/repositories/user"
	"net/http"
)

type AuthController struct {
	db *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{db: db}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (c *AuthController) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if authRequest.Username == "" || authRequest.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	userRepo := user.NewRepository(c.db)
	user, err := userRepo.FindByUsernameAndPassword(authRequest.Username, authRequest.Password)
	if err != nil {
		log.Printf("Ошибка в базе данных: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user == nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	token, err := generateToken()
	if err != nil {
		log.Println("Ошибка при создании токена", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenDurationHours := 24
	authResponse := auth_token.NewRepository(c.db)
	_, err = authResponse.FindAll(1, token, tokenDurationHours)
	if err != nil {
		log.Println("Ошибка сохранения токена", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
