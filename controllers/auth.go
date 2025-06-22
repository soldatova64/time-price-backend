package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"main/entity"
	"main/models/request"
	"main/models/responses"
	"main/repositories/auth_token"
	"main/repositories/user"
	"net/http"
)

func (c *App) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var authRequest request.AuthRequest
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
	authTokenRepo := auth_token.NewRepository(c.db)

	authToken := &entity.AuthToken{
		UserID: user.ID,
		Token:  token,
	}

	_, err = authTokenRepo.AddToken(authToken, tokenDurationHours)
	if err != nil {
		log.Println("Token storage error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := responses.AuthResponse{
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
