package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"net/http"
)

func main() {
	protocol, err := cloudevents.NewHTTP()

	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}
	h, err := cloudevents.NewHTTPReceiveHandler(context.Background(), protocol, receive)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	log.Printf("will listen on :8081\n")
	go func() {
		if err := http.ListenAndServe(":8081", h); err != nil {
			log.Fatalf("unable to start http server, %s", err)
		}
	}()

	log.Print("Starting Result Service")

	mux := http.NewServeMux()
	mux.Handle("/", middleware())
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("unable to start http server for Result Service, %s", err)
	}
}

func receive(event cloudevents.Event) {
	log.Printf("%s", event)
}

func middleware() http.Handler {
	return &handler{}
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Get was triggered"))

}
