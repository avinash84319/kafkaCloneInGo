package handlers

import (
	"fmt"
	"net/http"
)

func HeadlthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome To Kafka Clone In GO")
}