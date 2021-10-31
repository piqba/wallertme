package app

import (
	"context"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Run(ctx context.Context) error {
	var span trace.Span
	_, span = otel.Tracer("example").Start(ctx, "Run")
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		span.End()
		return err
	}

	address := "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	_, err = cardano.InfoByAddress(context.TODO(), address)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(sumary.Result.ToJSON())
	span.End()
	return nil
}
