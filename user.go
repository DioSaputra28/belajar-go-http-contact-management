package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Users struct {
	UserId    int64   `json:"user_id"`
	Name      string  `json:"name" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password,omitempty" validate:"required"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
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

func UserLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	
	if err := validate.StructPartial(user, "Email", "Password"); err != nil {
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

	var userData Users
	sha1 := sha1.New()
	sha1.Write([]byte(user.Password))
	hashedPassword := fmt.Sprintf("%x", sha1.Sum(nil))

	err = db.QueryRow("SELECT user_id, name, email, created_at, updated_at FROM users WHERE email = ? AND password = ?", user.Email, hashedPassword).Scan(&userData.UserId, &userData.Name, &userData.Email, &userData.CreatedAt, &userData.UpdatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Invalid email or password",
		})
		return
	}

	token := uuid.New().String()
	_, err = db.Exec("UPDATE users SET token = ? WHERE user_id = ?", token, userData.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	userDataMap := map[string]any{
		"user_id":    userData.UserId,
		"email":      userData.Email,
		"token":      token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Login successful",
		"user":    userDataMap,
	})
}
func GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := Connect()
	if err != nil {
		fmt.Println("Error connect:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		fmt.Println("Error count:", err)
	} else {
		fmt.Println("Jumlah data di database:", count)
	}

	data, err := db.Query("SELECT user_id, name, email, created_at, updated_at FROM users")
	if err != nil {
		fmt.Println("Error query:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}
	defer data.Close()

	users := []Users{}

	for data.Next() {
		var user Users
		err := data.Scan(&user.UserId, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println("Error scan:", err)
			continue
		}
		users = append(users, user)
	}

	if err := data.Err(); err != nil {
		fmt.Println("Error rows:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
		"data":    users,
	})
}

func GetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := Connect()
	if err != nil {
		fmt.Println("Error connect:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}
	defer db.Close()

	var user Users
	err = db.QueryRow("SELECT user_id, name, email, created_at, updated_at FROM users WHERE user_id = ?", ps.ByName("id")).Scan(&user.UserId, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Println("Error query:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "User not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
		"data":    user,
	})
}

func UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var userId Users
	err = db.QueryRow("SELECT user_id, name, email, created_at, updated_at FROM users WHERE user_id = ?", ps.ByName("id")).Scan(&userId.UserId, &userId.Name, &userId.Email, &userId.CreatedAt, &userId.UpdatedAt)
	if err != nil {
		fmt.Println("Error query:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "User not found",
		})
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE user_id = ?", user.Name, user.Email, ps.ByName("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User updated successfully",
		"user":    userId,
	})
}
