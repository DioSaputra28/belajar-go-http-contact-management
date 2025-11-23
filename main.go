package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
    router := httprouter.New()
    router.POST("/user", CreateUser)

    log.Println("Server running on http://localhost:8080") 
    log.Fatal(http.ListenAndServe(":8080", router))
}
