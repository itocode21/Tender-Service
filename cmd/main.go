package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/itocode21/Tender-Service/internal/handlers"
	"github.com/itocode21/Tender-Service/internal/repository"
	"github.com/itocode21/Tender-Service/internal/routers"
	"github.com/itocode21/Tender-Service/internal/services"
	postgresqldb "github.com/itocode21/Tender-Service/postgresql-db"
)

func main() {
	godotenv.Load()
	db, err := postgresqldb.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer db.Close()

	organizationRepo := repository.NewOrganizationRepository(db)
	//employeeRepo := repository.NewEmployeeRepository(db)
	tenderRepo := repository.NewTenderRepository(db)
	proposalRepo := repository.NewProposalRepository(db)

	organizationService := services.NewOrganizationService(organizationRepo)
	tenderService := services.NewTenderService(tenderRepo)
	proposalService := services.NewProposalService(proposalRepo, tenderRepo)

	organizationHandler := handlers.NewOrganizationHandler(organizationService)
	tenderHandler := handlers.NewTenderHandler(tenderService)
	proposalHandler := handlers.NewProposalHandler(proposalService)

	router := routers.SetupRoutes(organizationHandler, tenderHandler, proposalHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
