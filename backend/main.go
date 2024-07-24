package main

import (
	"fmt"
	"jwtAuth/router"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	fmt.Println("JWT Auth ..........")
    r := router.Router()
    headers := handlers.AllowedHeaders([]string{"Content-Type","Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methods := handlers.AllowedMethods([]string{"GET","POST","PUT","DELETE","OPTIONS"})

	log.Fatal(http.ListenAndServe(":8080",handlers.CORS(headers,origins,methods)(r)))
	fmt.Println("Server is running on Port 8080 ..........")

}