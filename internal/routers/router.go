package routers

//сначала закинул в мейн но как то  мусорно вышло, поэтому будет тут, не зря же практика хорошая

import (
	"github.com/gorilla/mux"
	"github.com/itocode21/Tender-Service/internal/handlers"
)

func SetupRoutes(organizationHandler *handlers.OrganizationHandler, tenderHandler *handlers.TenderHandler, proposalHandler *handlers.ProposalHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/organizations", organizationHandler.CreateOrganizationHandler).Methods("POST")
	router.HandleFunc("/organizations", organizationHandler.ListOrganizationsHandler).Methods("GET")
	router.HandleFunc("/organizations/{id}", organizationHandler.GetOrganizationHandler).Methods("GET")
	router.HandleFunc("/organizations/{id}", organizationHandler.UpdateOrganizationHandler).Methods("PUT")
	router.HandleFunc("/organizations/{id}", organizationHandler.DeleteOrganizationHandler).Methods("DELETE")
	router.HandleFunc("/organizations/{id}/responsibles", organizationHandler.AddResponsibleHandler).Methods("POST")
	router.HandleFunc("/organizations/{id}/responsibles", organizationHandler.RemoveResponsibleHandler).Methods("DELETE")

	router.HandleFunc("/tenders", tenderHandler.CreateTenderHandler).Methods("POST")
	router.HandleFunc("/tenders", tenderHandler.ListTendersHandler).Methods("GET")
	router.HandleFunc("/tenders/{id}", tenderHandler.GetTenderHandler).Methods("GET")
	router.HandleFunc("/tenders/{id}", tenderHandler.UpdateTenderHandler).Methods("PUT")
	router.HandleFunc("/tenders/{id}", tenderHandler.DeleteTenderHandler).Methods("DELETE")
	router.HandleFunc("/tenders/organization/{organization_id}", tenderHandler.GetTendersByOrganizationHandler).Methods("GET")
	router.HandleFunc("/tenders/{id}/publish", tenderHandler.PublishTenderHandler).Methods("PUT")
	router.HandleFunc("/tenders/{id}/close", tenderHandler.CloseTenderHandler).Methods("PUT")
	router.HandleFunc("/tenders/{id}/cancel", tenderHandler.CancelTenderHandler).Methods("PUT")

	router.HandleFunc("/proposals", proposalHandler.CreateProposalHandler).Methods("POST")
	router.HandleFunc("/proposals/{id}", proposalHandler.GetProposalHandler).Methods("GET")
	router.HandleFunc("/proposals/{id}", proposalHandler.UpdateProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}", proposalHandler.DeleteProposalHandler).Methods("DELETE")
	router.HandleFunc("/proposals", proposalHandler.ListProposalsHandler).Methods("GET")
	router.HandleFunc("/tenders/{tender_id}/proposals", proposalHandler.GetProposalsByTenderHandler).Methods("GET")
	router.HandleFunc("/proposals/{id}/publish", proposalHandler.PublishProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/accept", proposalHandler.AcceptProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/reject", proposalHandler.RejectProposalHandler).Methods("PUT")
	router.HandleFunc("/proposals/{id}/cancel", proposalHandler.CancelProposalHandler).Methods("PUT")

	return router
}
