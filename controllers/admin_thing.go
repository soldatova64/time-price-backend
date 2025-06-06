package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/entity"
	"main/models"
	"main/models/responses"
	"main/repositories/thing"
	"net/http"
	"strconv"
	"time"
)

func (app *App) AdminThingController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_thing"}

	var thingEntity entity.Thing

	if err := json.NewDecoder(request.Body).Decode(&thingEntity); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "json",
					Message: "Недопустимый формат JSON",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var errors []responses.Error

	if thingEntity.Name == "" {
		errors = append(errors, responses.Error{
			Field:   "name",
			Message: "Поле должно быть заполнено",
		})
	}

	if thingEntity.PayDate.IsZero() {
		errors = append(errors, responses.Error{
			Field:   "pay_date",
			Message: "Необходимо указать дату покупки",
		})
	} else if thingEntity.PayDate.After(time.Now()) {
		errors = append(errors, responses.Error{
			Field:   "pay_date",
			Message: "Дата не может быть в будущем",
		})
	}

	if thingEntity.PayPrice <= 0 {
		errors = append(errors, responses.Error{
			Field:   "pay_price",
			Message: "Стоимость покупки должна быть положительной",
		})
	}

	if len(errors) > 0 {
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	thingRepo := thing.NewRepository(app.db)
	createdThing, err := thingRepo.Add(app.db, &thingEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось добавить вещь в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	successResponse := responses.AdminThingResponse{
		Meta: meta,
		Data: *createdThing,
	}
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}

func (app *App) AdminThingUpdateController(writer http.ResponseWriter, request *http.Request) {
	meta := models.Meta{Action: "admin_thing_update"}

	vars := mux.Vars(request)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "id",
					Message: "Неверный формат ID",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	var thingEntity entity.Thing
	if err := json.NewDecoder(request.Body).Decode(&thingEntity); err != nil {
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "json",
					Message: "Недопустимый формат JSON",
				},
			},
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	thingEntity.ID = id
	var errors []responses.Error
	if thingEntity.Name == "" {
		errors = append(errors, responses.Error{
			Field:   "name",
			Message: "Поле должно быть заполнено",
		})
	}
	if thingEntity.PayDate.IsZero() {
		errors = append(errors, responses.Error{
			Field:   "pay_date",
			Message: "Необходимо указать дату покупки",
		})
	} else if thingEntity.PayDate.After(time.Now()) {
		errors = append(errors, responses.Error{
			Field:   "pay_date",
			Message: "Дата не может быть в будущем",
		})
	}
	if thingEntity.PayPrice <= 0 {
		errors = append(errors, responses.Error{
			Field:   "pay_price",
			Message: "Стоимость покупки должна быть положительной",
		})
	}
	if len(errors) > 0 {
		errorResponse := responses.ErrorResponse{
			Meta:   meta,
			Errors: errors,
		}
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}

	thingRepo := thing.NewRepository(app.db)
	updatedThing, err := thingRepo.Update(app.db, &thingEntity)
	if err != nil {
		log.Printf("Database error: %v", err)
		errorResponse := responses.ErrorResponse{
			Meta: meta,
			Errors: []responses.Error{
				{
					Field:   "database",
					Message: "Не удалось обновить вещь в БД",
				},
			},
		}
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(errorResponse)
		return
	}
	successResponse := responses.AdminThingResponse{
		Meta: meta,
		Data: *updatedThing,
	}
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
