package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"zadanie-6105/pkg/utils"
)

type contextKey string

const userContextKey = contextKey("username")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем поле "creatorUsername" в теле запроса
		var body struct {
			CreatorUsername string `json:"creatorUsername"`
		}

		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Println("Failed to parse request body")
				utils.RespondWithError(w, http.StatusBadRequest, "Неверный формат запроса")
				return
			}
		} else {
			body.CreatorUsername = r.URL.Query().Get("username")
		}

		if body.CreatorUsername == "" {
			log.Println("creatorUsername is missing")
			utils.RespondWithError(w, http.StatusUnauthorized, "Пользователь не аутентифицирован")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, body.CreatorUsername)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(userContextKey).(string)
	return username, ok
}
