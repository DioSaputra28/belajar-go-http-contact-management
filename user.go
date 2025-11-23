package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Users struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	var user Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		errMsgs := []string{}
		for _, e := range errors {
			errMsgs = append(errMsgs, fmt.Sprintf("%s is %s", e.Field(), e.ActualTag()))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": errMsgs,
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

	data := &user
	sha1 := sha1.New()
	sha1.Write([]byte(data.Password))
	data.Password = fmt.Sprintf("%x", sha1.Sum(nil))

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", data.Name, data.Email, data.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User created successfully",
		"user":    user,
	})
}
