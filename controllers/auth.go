package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"main/entity"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/auth_token"
	"main/repositories/user"
	"net/http"
)

func (c *App) AuthHandler(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "auth_token"}

	if request.Method != http.MethodPost {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "method",
					Message: "Разрешен только метод POST",
				},
			},
		}
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var authRequest requests.AuthRequest
	if err := json.NewDecoder(request.Body).Decode(&authRequest); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "request",
					Message: "Недопустимый формат запроса",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if authRequest.Username == "" || authRequest.Password == "" {
		errors := []responses.Error{}
		if authRequest.Username == "" {
			errors = append(errors, responses.Error{
				Field:   "username",
				Message: "Требуется имя пользователя",
			})
		}
		if authRequest.Password == "" {
			errors = append(errors, responses.Error{
				Field:   "password",
				Message: "Требуется ввести пароль",
			})
		}

		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	userRepo := user.NewRepository(c.db)
	user, err := userRepo.FindByUsernameAndPassword(authRequest.Username, authRequest.Password)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Не удалось выполнить проверку подлинности",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if user == nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "credentials",
					Message: "Неверное имя пользователя или пароль",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	token, err := generateToken()
	if err != nil {
		log.Println("Token generation error:", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Внутренняя ошибка сервера",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
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
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Внутренняя ошибка сервера",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	response := struct {
		Meta models.Meta            `json:"meta"`
		Data responses.AuthResponse `json:"data"`
	}{
		Meta: meta,
		Data: responses.AuthResponse{
			Token: token,
		},
	}
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
