package middleware

import (
	"context"
	"database/sql"
	"log"
	"main/repositories/auth_token"
	"net/http"
	"strings"
)

func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем OPTIONS запросы и публичные маршруты
			if r.Method == "OPTIONS" ||
				r.URL.Path == "/auth" ||
				r.URL.Path == "/register" {
				log.Printf("AuthMiddleware: skipping auth for %s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("AuthMiddleware: protecting %s %s", r.Method, r.URL.Path)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("AuthMiddleware: missing Authorization header")
				http.Error(w, `{"error": "Требуется авторизация"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Printf("AuthMiddleware: invalid Authorization header format")
				http.Error(w, `{"error": "Неверный формат токена"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(parts[1])
			log.Printf("AuthMiddleware: Received token: %s", token)

			authTokenRepo := auth_token.NewRepository(db)
			authToken, err := authTokenRepo.FindByToken(token)
			if err != nil {
				log.Printf("AuthMiddleware: Token not found or error: %v", err)
				http.Error(w, `{"error": "Неверный токен"}`, http.StatusUnauthorized)
				return
			}

			log.Printf("AuthMiddleware: User authenticated, user_id: %d", authToken.UserID)

			// Добавляем user_id в контекст
			ctx := context.WithValue(r.Context(), "user_id", authToken.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
