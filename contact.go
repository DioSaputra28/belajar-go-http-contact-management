package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Contacts struct {
	ContactId string  `json:"contact_id"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     string  `json:"phone" validate:"required"`
	UserId    string  `json:"user_id"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}

func CreateContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	var contact Contacts
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(contact); err != nil {
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

	ctxUser := r.Context().Value("user").(Users)

	_, err = db.Exec("INSERT INTO contacts (first_name, last_name, email, phone, user_id) VALUES (?, ?, ?, ?, ?)", contact.FirstName, contact.LastName, contact.Email, contact.Phone, ctxUser.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	data := map[string]any{
		"first_name": contact.FirstName,
		"last_name":  contact.LastName,
		"email":      contact.Email,
		"phone":      contact.Phone,
		"user_id":    ctxUser.UserId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Contact created successfully",
		"data": data,
	})

}

func GetContacts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	ctxUser := r.Context().Value("user").(Users)

	rows, err := db.Query("SELECT * FROM contacts WHERE user_id = ?", ctxUser.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}
	defer rows.Close()

	var contacts []Contacts
	for rows.Next() {
		var contact Contacts
		if err := rows.Scan(&contact.ContactId, &contact.FirstName, &contact.LastName, &contact.Email, &contact.Phone, &contact.UserId, &contact.CreatedAt, &contact.UpdatedAt); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Internal Server Error",
			})
			return
		}
		contacts = append(contacts, contact)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
		"data":    contacts,
	})
}

func GetContactId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	ctxUser := r.Context().Value("user").(Users)

	var contact Contacts

	err = db.QueryRow("SELECT * FROM contacts WHERE contact_id = ? AND user_id = ?", ps.ByName("id"), ctxUser.UserId).Scan(&contact.ContactId, &contact.FirstName, &contact.LastName, &contact.Email, &contact.Phone, &contact.UserId, &contact.CreatedAt, &contact.UpdatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Contact not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
		"data":    contact,
	})
}

func UpdateContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	var contact Contacts
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Input tidak valid.",
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(contact); err != nil {
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

	ctxUser := r.Context().Value("user").(Users)


	rows, err := db.Exec("UPDATE contacts SET first_name = ?, last_name = ?, email = ?, phone = ? WHERE contact_id = ? AND user_id = ?", contact.FirstName, contact.LastName, contact.Email, contact.Phone, ps.ByName("id"), ctxUser.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	rowsAffected, _ :=  rows.RowsAffected()
	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Contact not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Contact updated successfully",
		"data": map[string]any{
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"email":      contact.Email,
			"phone":      contact.Phone,
		},
	})
}

func DeleteContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	ctxUser := r.Context().Value("user").(Users)

	rows, err := db.Exec("DELETE FROM contacts WHERE contact_id = ? AND user_id = ?", ps.ByName("id"), ctxUser.UserId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Internal Server Error",
		})
		return
	}

	rowsAffected, _ :=  rows.RowsAffected()
	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Contact not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Contact deleted successfully",
	})
}