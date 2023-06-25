//go:generate weaver generate ./...
package main

import (
	"context"
	"log"

	"github.com/ServiceWeaver/weaver"
	"github.com/syllabix/kafkaless/consumer"
	"github.com/syllabix/kafkaless/web"
)

type app struct {
	weaver.Implements[weaver.Main]

	// application service dependencies
	_ weaver.Ref[web.Server]
	_ weaver.Ref[consumer.Service]
}

func boot(ctx context.Context, app *app) error {
	// block until application context is done
	<-ctx.Done()
	// in an ideal world we could potentially configure
	// graceful termination here - but service weaver invokes
	// os.Exit on SIGINT and SIGTERM which will kill the process
	// immediately.
	// there is an open discussion about adding a shutdown
	// lifecycle hook (similar to component.Init(...)) being discussed here
	// https://github.com/ServiceWeaver/weaver/issues/275
	return ctx.Err()
}

func main() {
	if err := weaver.Run(context.Background(), boot); err != nil {
		log.Fatal(err)
	}
}
