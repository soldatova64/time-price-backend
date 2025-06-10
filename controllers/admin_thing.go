package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/entity"
	"main/models"
	"main/models/responses"
	"main/repositories/thing"
	"main/types"
	"net/http"
	"strconv"
	"strings"
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

	path := strings.TrimPrefix(request.URL.Path, "/admin/thing/")
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
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

	var updateData map[string]interface{}
	if err := json.NewDecoder(request.Body).Decode(&updateData); err != nil {
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

	thingEntity := entity.Thing{ID: id}

	if name, ok := updateData["name"].(string); ok {
		thingEntity.Name = name
	}

	if payDateStr, ok := updateData["pay_date"].(string); ok {
		payDate, err := time.Parse(time.RFC3339, payDateStr)
		if err != nil {
			errorResponse := responses.ErrorResponse{
				Meta: meta,
				Errors: []responses.Error{
					{
						Field:   "pay_date",
						Message: "Неверный формат даты",
					},
				},
			}
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(errorResponse)
			return
		}
		thingEntity.PayDate = payDate
	}

	if payPrice, ok := updateData["pay_price"].(float64); ok {
		thingEntity.PayPrice = int(payPrice)
	}

	if saleDate, ok := updateData["sale_date"]; ok {
		if saleDate == nil {
			thingEntity.SaleDate = types.NullTime{NullTime: sql.NullTime{Valid: false}}
		} else {
			parsedDate, err := time.Parse(time.RFC3339, saleDate.(string))
			if err != nil {
				errorResponse := responses.ErrorResponse{
					Meta: meta,
					Errors: []responses.Error{
						{
							Field:   "sale_date",
							Message: "Неверный формат даты",
						},
					},
				}
				writer.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(writer).Encode(errorResponse)
				return
			}
			thingEntity.SaleDate = types.NullTime{NullTime: sql.NullTime{
				Time:  parsedDate,
				Valid: true,
			}}
		}
	}

	if salePrice, ok := updateData["sale_price"]; ok {
		if salePrice == nil {
			thingEntity.SalePrice = types.NullInt64{NullInt64: sql.NullInt64{Valid: false}}
		} else {
			thingEntity.SalePrice = types.NullInt64{NullInt64: sql.NullInt64{
				Int64: int64(salePrice.(float64)),
				Valid: true,
			}}
		}
	}

	var errors []responses.Error

	if val, ok := updateData["name"]; ok && val == "" {
		errors = append(errors, responses.Error{
			Field:   "name",
			Message: "Имя не может быть пустым",
		})
	}

	if _, ok := updateData["pay_date"]; ok {
		if thingEntity.PayDate.IsZero() {
			errors = append(errors, responses.Error{
				Field:   "pay_date",
				Message: "Неверный формат даты",
			})
		} else if thingEntity.PayDate.After(time.Now()) {
			errors = append(errors, responses.Error{
				Field:   "pay_date",
				Message: "Дата не может быть в будущем",
			})
		}
	}

	if val, ok := updateData["pay_price"]; ok && val.(float64) <= 0 {
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
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(successResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(writer, "Ошибка формирования JSON", http.StatusInternalServerError)
	}
}
