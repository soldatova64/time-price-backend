package controllers

import (
	"encoding/json"
	"log"
	"main/entity"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/user"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func init() {
	var validate = validator.New()
	validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
			return false
		}
		if email[0] == '@' || email[len(email)-1] == '@' {
			return false
		}

		parts := strings.Split(email, ".")
		if len(parts) < 2 || len(parts[len(parts)-1]) < 2 {
			return false
		}

		return true
	})
}

func (app *App) AdminUserController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_user"}

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
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var req requests.UserRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
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

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		errors := parseValidationErrors(err)
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errors)
		return
	}

	passwordHash := hashPassword(req.Password)

	userRepo := user.NewRepository(app.db)
	newUser := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: passwordHash,
	}

	newUser, err := userRepo.Add(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Не удалось создать пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	response := struct {
		Meta models.Meta `json:"meta"`
		Data struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}{
		Meta: meta,
		Data: struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       newUser.ID,
			Username: newUser.Username,
			Email:    newUser.Email,
		},
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func parseValidationErrors(err error) []responses.Error {
	var errors []responses.Error

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		var message string

		switch err.Tag() {
		case "required":
			message = "Поле обязательно для заполнения"
		case "min":
			if field == "Username" {
				message = "Имя пользователя должно содержать минимум 3 символа"
			} else {
				message = "Пароль должен содержать минимум 6 символов"
			}
		case "email":
			message = "Некорректный формат email"
		default:
			message = "Недопустимое значение"
		}

		errors = append(errors, responses.Error{
			Field:   field,
			Message: message,
		})
	}

	return errors
}

func hashPassword(password string) string {
	return "hashed_" + password
}
