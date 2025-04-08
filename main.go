package main

import (
	"log"
	"net/http"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/routes"
	_ "github.com/Anwarjondev/task-management-api/docs" // Import generated docs
    httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db.Connect()
	db.AutoMigrate()
	mux := routes.SetUpRoutes()


	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running on port :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server Failed: ", err)
	}
}