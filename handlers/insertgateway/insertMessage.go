package insertgateway

import (
	"encoding/json"
	"net/http"

	"github.com/avinash84319/kafkaCloneInGo/models/insertgateway"
	"github.com/avinash84319/kafkaCloneInGo/systems"
)

func ReciveMessage(w http.ResponseWriter, r *http.Request) {

	var newRequest insertgateway.Request
	
	//parse the request
	_ = json.NewDecoder(r.Body).Decode(&newRequest)

	//submit the request to topic handler goroutine
	systems.MasterTopicFunction(newRequest)

	// send response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received successfully"))
}