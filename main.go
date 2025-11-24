package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
    router := httprouter.New()
    router.POST("/user", CreateUser)
    router.POST("/login", UserLogin)
    router.GET("/user", AuthMiddleware(GetUser))
    router.GET("/user/:id", AuthMiddleware(GetUserId))
    router.PUT("/user/:id", AuthMiddleware(UpdateUser))

    log.Println("Server running on http://localhost:8080") 
    log.Fatal(http.ListenAndServe(":8080", router))
}
