package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/handlers"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
)

func main() {

	db, err := postgresqldb.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/api/users", employeeHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}", employeeHandler.GetUserHandler).Methods("GET")

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", r))

}
