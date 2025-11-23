package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:Diosaputra288!@/contact_management")
	if err != nil {
		return nil, err
	}
	return db, nil
}
