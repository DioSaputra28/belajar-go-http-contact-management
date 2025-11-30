package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB - Initialize database connection pool (call once at startup)
func InitDB() error {
	var err error
	once.Do(func() {
		// Load .env hanya sekali
		_ = godotenv.Load()

		dbUser := getEnv("DB_USER", "root")
		dbPassword := getEnv("DB_PASSWORD", "")
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "3306")
		dbName := getEnv("DB_NAME", "contact_management")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbUser, dbPassword, dbHost, dbPort, dbName)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return
		}

		db.SetMaxIdleConns(50)  
		db.SetMaxOpenConns(200) 
		db.SetConnMaxIdleTime(10 * time.Minute)
		db.SetConnMaxLifetime(30 * time.Minute)

		// Test connection
		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v", err)
			return
		}

		log.Println("Database connection pool initialized successfully")
	})
	return err
}

// GetDB - Get database connection pool
func GetDB() *sql.DB {
	return db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
