//go:generate weaver generate ./...
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/syllabix/kafkaless/reverser"
)

func main() {
	if err := weaver.Run(context.Background(), start); err != nil {
		log.Fatal(err)
	}
}

// app is the main component of the application. weaver.Run creates
// it and passes it to serve.
type app struct {
	weaver.Implements[weaver.Main]
	reverser weaver.Ref[reverser.Service]
	server   weaver.Listener
}

func start(ctx context.Context, app *app) error {
	// The server listener will listen on a random port chosen by the operating
	// system. This behavior can be changed in the config file.
	fmt.Printf("server listener available on %v\n", app.server)

	// Serve the /hello endpoint.
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		reversed, err := app.reverser.Get().Reverse(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hello, %s!\n", reversed)
	})

	return http.Serve(app.server, nil)
}
