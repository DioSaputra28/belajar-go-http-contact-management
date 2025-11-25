package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Addresses struct {
	AddressId  string  `json:"address_id"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	Province   string  `json:"province"`
	Country    string  `json:"country" validate:"required"`
	PostalCode string  `json:"postal_code"`
	ContactId  string  `json:"contact_id" validate:"required"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
}

func CreateAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Request body is empty",
		})
	}

	var address Addresses
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	validate := validator.New()
	if err = validate.Struct(address); err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM contacts WHERE contact_id = ?", address.ContactId).Scan(&count)
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contact not found",
		})
		return
	}

	_, err = db.Exec("INSERT INTO addresses (street, city, province, country, postal_code, contact_id) VALUES (?, ?, ?, ?, ?, ?)", address.Street, address.City, address.Province, address.Country, address.PostalCode, address.ContactId)
	if err != nil {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	data := map[string]any{
		"street":      address.Street,
		"city":        address.City,
		"province":    address.Province,
		"country":     address.Country,
		"postal_code": address.PostalCode,
		"contact_id":  address.ContactId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Address created successfully",
		"data":    data,
	})
}

func GetAddresses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := Connect()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM contacts WHERE contact_id = ?", ps.ByName("contactId")).Scan(&count)
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contact not found",
		})
		return
	}

	rows, err := db.Query("SELECT * FROM addresses WHERE contact_id = ?", ps.ByName("contactId"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer rows.Close()

	var addresses []Addresses
	for rows.Next() {
		var address Addresses
		err = rows.Scan(&address.AddressId, &address.Street, &address.City, &address.Province, &address.Country, &address.PostalCode, &address.ContactId, &address.CreatedAt, &address.UpdatedAt)
		if err != nil {
			fmt.Println("Error disini", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Internal server error",
			})
			return
		}
		addresses = append(addresses, address)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Addresses retrieved successfully",
		"data":    addresses,
	})
}

func GetAddressId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := Connect()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM contacts WHERE contact_id = ?", ps.ByName("contactId")).Scan(&count)
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contact not found",
		})
		return
	}

	var address Addresses
	err = db.QueryRow("SELECT * FROM addresses WHERE address_id = ?", ps.ByName("addressId")).Scan(&address.AddressId, &address.Street, &address.City, &address.Province, &address.Country, &address.PostalCode, &address.ContactId, &address.CreatedAt, &address.UpdatedAt)
	if err != nil {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Address Not Found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Address retrieved successfully",
		"data":    address,
	})
}

func UpdateAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("contactId:", ps.ByName("contactId"))
    fmt.Println("addressId:", ps.ByName("addressId"))

	if r.Body == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Request body is empty",
		})
	}

	var address Addresses
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request body",
		})
		return
	}

	validate := validator.New()
	if err = validate.StructPartial(address, "Country"); err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM contacts WHERE contact_id = ?", ps.ByName("contactId")).Scan(&count)
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contact not found",
		})
		return
	}

	row, err := db.Exec("UPDATE addresses SET street = ?, city = ?, province = ?, country = ?, postal_code = ? WHERE address_id = ? AND contact_id = ?", address.Street, address.City, address.Province, address.Country, address.PostalCode, ps.ByName("addressId"), ps.ByName("contactId"))
	if err != nil {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	countRow, err := row.RowsAffected()
	if countRow == 0 {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Address not found",
		})
		return
	}

	data := map[string]any{
		"street":      address.Street,
		"city":        address.City,
		"province":    address.Province,
		"country":     address.Country,
		"postal_code": address.PostalCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Address updated successfully",
		"data":    data,
	})
}

func DeleteAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, err := Connect()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	defer db.Close()

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM contacts WHERE contact_id = ?", ps.ByName("contactId")).Scan(&count)
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contact not found",
		})
		return
	}

	row, err := db.Exec("DELETE FROM addresses WHERE address_id = ? AND contact_id = ?", ps.ByName("addressId"), ps.ByName("contactId"))
	if err != nil {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Internal server error",
		})
		return
	}

	countRow, err := row.RowsAffected()
	if countRow == 0 {
		fmt.Println("Error disini", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Address not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Address deleted successfully",
	})
}
