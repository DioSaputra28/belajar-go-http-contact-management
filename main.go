package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	_ = godotenv.Load()

	// Initialize database connection pool
	if err := InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer GetDB().Close()

	log.Println("Starting Contact Management API...")

	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/login", UserLogin)
	router.GET("/user", AuthMiddleware(GetUser))
	router.GET("/user/:id", AuthMiddleware(GetUserId))
	router.PUT("/user/:id", AuthMiddleware(UpdateUser))

	router.POST("/contact", AuthMiddleware(CreateContact))
	router.GET("/contact", AuthMiddleware(GetContacts))
	router.GET("/contact/:id", AuthMiddleware(GetContactId))
	router.PUT("/contact/:id", AuthMiddleware(UpdateContact))
	router.DELETE("/contact/:id", AuthMiddleware(DeleteContact))

	router.POST("/address/", AuthMiddleware(CreateAddress))
	router.GET("/address/:contactId", AuthMiddleware(GetAddresses))
	router.GET("/address/:contactId/:addressId", AuthMiddleware(GetAddressId))
	router.PUT("/address/:contactId/:addressId", AuthMiddleware(UpdateAddress))
	router.DELETE("/address/:contactId/:addressId", AuthMiddleware(DeleteAddress))

	router.GET("/docs/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.ServeFile(w, r, "docs"+ps.ByName("filepath"))
	})

	router.GET("/swagger/*any", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		httpSwagger.Handler(
			httpSwagger.URL("/docs/swagger.yaml"),
		).ServeHTTP(w, r)
	})

	port := getEnv("APP_PORT", "8080")

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
