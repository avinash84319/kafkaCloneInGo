package route

import (
	"github.com/gorilla/mux"

	"github.com/avinash84319/kafkaCloneInGo/handlers/insertgateway"
)

func AddInsertGatewayRoutes(router *mux.Router) {
	// recieve message route
	router.HandleFunc("/sendMessage", insertgateway.ReciveMessage).Methods("POST")
}
