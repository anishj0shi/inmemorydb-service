package main

import (
	"github.com/anishj0shi/inmemorydb-service/pkg/client"
	"log"
	"net/http"
)

var eventHandler = client.NewEventResultSrvice()

func main() {
	http.HandleFunc("/eventResult", HandleAddResult)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("unable to start http server for Result Service, %s", err)
	}
}

func HandleAddResult(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		eventHandler.GetEventResult(w, req)
	case http.MethodPost:
		eventHandler.PostEventResult(w, req)

	}
}
