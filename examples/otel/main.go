package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/piqba/wallertme/examples/otel/app"
	"github.com/piqba/wallertme/pkg/otelify"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	l := log.New(os.Stdout, "", 0)
	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	expo, err := otelify.NewExporter(f)
	if err != nil {
		l.Fatal(err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(expo),
		trace.WithResource(
			otelify.NewResource(
				"example",
				"v0.3.2",
				"dev",
			),
		),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	errCh := make(chan error)
	go func() {
		errCh <- app.Run(context.Background())
	}()
	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
