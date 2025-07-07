package helpers

import (
	"github.com/go-playground/validator/v10"
	"main/models/responses"
	"regexp"
	"time"
)

func SimpleEmailValidation(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func ParseValidationErrors(err error) []responses.Error {
	var errors []responses.Error

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		var message string

		switch err.Tag() {
		case "required":
			message = "Поле обязательно для заполнения"
		case "min":
			if field == "Name" || field == "Username" || field == "Description" {
				message = "Поле должно содержать минимум 3 символа"
			} else {
				message = "Поле должно содержать минимум " + err.Param() + " символов"
			}
		case "email":
			message = "Некорректный формат email"
		case "gt":
			if field == "ThingID" {
				message = "Необходимо указать корректный ID вещи"
			} else {
				message = "Значение должно быть больше 0"
			}
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

func IsFutureDate(date time.Time) bool {
	return date.After(time.Now())
}

func IsPastDate(date time.Time) bool {
	return date.Before(time.Now())
}
