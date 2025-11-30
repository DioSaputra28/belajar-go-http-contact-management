package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		token := r.Header.Get("Authorization")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Unauthorized",
			})
			return
		}

		db := GetDB()

		var user Users
		err := db.QueryRow("SELECT user_id, name, email, created_at FROM users WHERE token = ?", token).Scan(&user.UserId, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Unauthorized",
			})
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next(w, r, p)
	}
}
