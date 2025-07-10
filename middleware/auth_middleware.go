package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"main/models/responses"
	"main/repositories/auth_token"
	"net/http"
	"time"
)

func AuthMiddleware(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			log.Printf("Получен токен: %s", token)
			if token == "" {
				log.Printf("Токен отсутствует")
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			repo := auth_token.NewRepository(db)
			authToken, err := repo.FindByToken(token)
			if err != nil {
				log.Printf("Ошибка поиска токена: %v", err)
				json.NewEncoder(w).Encode(responses.ErrorResponse{
					Errors: []responses.Error{
						{
							Field:   "auth",
							Message: "Ошибка проверки токена",
						},
					},
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if authToken.ID == 0 {
				log.Printf("Токен не найден в БД")
				json.NewEncoder(w).Encode(responses.ErrorResponse{
					Errors: []responses.Error{
						{
							Field:   "auth",
							Message: "Неверный или просроченный токен",
						},
					},
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if authToken.EndDate.Before(time.Now()) {
				json.NewEncoder(w).Encode(responses.ErrorResponse{
					Errors: []responses.Error{
						{
							Field:   "auth",
							Message: "Срок действия токена истек",
						},
					},
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", authToken.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
