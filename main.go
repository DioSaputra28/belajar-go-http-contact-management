package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
    router := httprouter.New()
    router.POST("/user", CreateUser)
    router.GET("/user", GetUser)
    router.GET("/user/:id", GetUserId)

    log.Println("Server running on http://localhost:8080") 
    log.Fatal(http.ListenAndServe(":8080", router))
}
