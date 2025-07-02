package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"main/entity"
	"main/helpers"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/user"
	"net/http"
)

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

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		errors := parseValidationErrors(err)
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if !helpers.SimpleEmailValidation(req.Email) {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "email",
					Message: "Некорректный формат email",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	passwordHash := helpers.HashPassword(req.Password)

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
