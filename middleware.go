package main

import (
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

		db, err := Connect()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Internal Server Error",
			})
			return
		}
		defer db.Close()

		var userCount int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE token = ?", token).Scan(&userCount)
		if err != nil || userCount == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Unauthorized",
			})
			return
		}
		next(w, r, p)
	}
}
