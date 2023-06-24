//go:generate weaver generate ./...
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/syllabix/kafkaless/consumer"
	"github.com/syllabix/kafkaless/producer"
)

// app is the main component of the application. weaver.Run creates
// it and passes it to serve.
type app struct {
	weaver.Implements[weaver.Main]
	weaver.Ref[consumer.Service]
	producer weaver.Ref[producer.Service]
	server   weaver.Listener
}

func main() {
	if err := weaver.Run(context.Background(), start); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context, app *app) error {
	// The server listener will listen on a random port chosen by the operating
	// system. This behavior can be changed in the config file.
	fmt.Printf("server listener available on %v\n", app.server)

	// Serve the /emit endpoint.
	http.HandleFunc("/emit", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("event")
		if name == "" {
			name = "World"
		}

		err := app.producer.Get().EmitEvent(ctx, name)
		if err != nil {
			http.Error(w, "oops... sorry about that", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.Serve(app.server, nil)
}
