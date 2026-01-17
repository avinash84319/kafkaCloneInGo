package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"github.com/avinash84319/kafkaCloneInGo/route"
)

func main() {
	router := mux.NewRouter()
	route.AddHealthRoute(router)
	route.AddInsertGatewayRoutes(router)

	// start the server
	http.ListenAndServe(":8080", router)
}
