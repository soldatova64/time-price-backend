package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"main/entity"
	"main/helpers"
	"main/models"
	"main/models/requests"
	"main/models/responses"
	"main/repositories/user"
	"net/http"
	"strconv"
	"time"
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
		errors := helpers.ParseValidationErrors(err)
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

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Ошибка при обработке пароля",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	userRepo := user.NewRepository(app.db)
	newUser := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: passwordHash,
	}

	newUser, err = userRepo.Add(newUser)
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

func (app *App) RegisterController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "register"}

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
		errors := helpers.ParseValidationErrors(err)
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

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Ошибка при обработке пароля",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	userRepo := user.NewRepository(app.db)
	newUser := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: passwordHash,
	}

	newUser, err = userRepo.Add(newUser)
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

	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
func (app *App) AdminUserUpdateController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_user_update"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	userIDStr := vars["id"]

	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "id",
					Message: "Некорректный ID пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	currentUserID, ok := request.Context().Value("user_id").(int)
	if !ok {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Требуется авторизация",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if targetUserID != currentUserID {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "permission",
					Message: "У вас нет прав для обновления данных другого пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusForbidden)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var req requests.UserUpdateRequest
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
		errors := helpers.ParseValidationErrors(err)
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if req.Username == nil && req.Password == nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "request",
					Message: "Должно быть указано username или password",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	updateData := make(map[string]interface{})

	if req.Username != nil {
		updateData["username"] = *req.Username
	}

	if req.Password != nil {
		// Хешируем новый пароль
		passwordHash, err := helpers.HashPassword(*req.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "system",
						Message: "Ошибка при обработке пароля",
					},
				},
			}
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(errorResponse)
			return
		}
		updateData["password"] = passwordHash
	}

	userRepo := user.NewRepository(app.db)
	updatedUser, err := userRepo.Update(targetUserID, updateData)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "user",
						Message: "Пользователь не найден или был удален",
					},
				},
			}
			writer.WriteHeader(http.StatusNotFound)
			json.NewEncoder(writer).Encode(errorResponse)
		} else {
			log.Printf("Error updating user: %v", err)
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "system",
						Message: "Не удалось обновить пользователя",
					},
				},
			}
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorResponse)
		}
		return
	}

	response := struct {
		Meta models.Meta `json:"meta"`
		Data struct {
			ID        int       `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`
	}{
		Meta: meta,
		Data: struct {
			ID        int       `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			UpdatedAt time.Time `json:"updated_at"`
		}{
			ID:        updatedUser.ID,
			Username:  updatedUser.Username,
			Email:     updatedUser.Email,
			UpdatedAt: updatedUser.UpdatedAt,
		},
	}

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func (app *App) AdminUserDeleteController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_user_delete"}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	userIDStr := vars["id"]

	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "id",
					Message: "Некорректный ID пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	// Проверяем авторизацию
	currentUserID, ok := request.Context().Value("user_id").(int)
	if !ok {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "auth",
					Message: "Требуется авторизация",
				},
			},
		}
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	// Разрешаем удаление только своего собственного аккаунта
	if targetUserID != currentUserID {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "permission",
					Message: "Вы можете удалить только свой собственный аккаунт",
				},
			},
		}
		writer.WriteHeader(http.StatusForbidden)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	userRepo := user.NewRepository(app.db)

	// Сначала проверяем, существует ли пользователь
	existingUser, err := userRepo.FindByID(targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "user",
						Message: "Пользователь не найден",
					},
				},
			}
			writer.WriteHeader(http.StatusNotFound)
			json.NewEncoder(writer).Encode(errorResponse)
			return
		}

		log.Printf("Error finding user: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Ошибка при поиске пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	if existingUser.IsDeleted {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "user",
					Message: "Пользователь уже удален",
				},
			},
		}
		writer.WriteHeader(http.StatusConflict)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	err = userRepo.Delete(targetUserID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "system",
					Message: "Не удалось удалить пользователя",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	response := struct {
		Meta models.Meta `json:"meta"`
		Data struct {
			ID        int       `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Message   string    `json:"message"`
			DeletedAt time.Time `json:"deleted_at"`
		} `json:"data"`
	}{
		Meta: meta,
		Data: struct {
			ID        int       `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Message   string    `json:"message"`
			DeletedAt time.Time `json:"deleted_at"`
		}{
			ID:        targetUserID,
			Username:  existingUser.Username,
			Email:     existingUser.Email,
			Message:   "Ваш аккаунт был успешно удален",
			DeletedAt: time.Now(),
		},
	}

	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
