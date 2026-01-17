package route

import (
	"github.com/gorilla/mux"

	"github.com/avinash84319/kafkaCloneInGo/handlers"
)

func AddHealthRoute(router *mux.Router) {
	router.HandleFunc("/", handlers.HeadlthCheckHandler).Methods("GET")
}
