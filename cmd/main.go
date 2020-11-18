package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go"
	"log"
)

func main()  {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), receive));
}

func receive(event cloudevents.Event) {
	log.Printf("%s", event)
}