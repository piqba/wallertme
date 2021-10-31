package bb8

import (
	"context"
	"time"
)

const (
	// nameBb8 is the Tracer nameBb8 used to identify this instrumentation library.
	nameBb8 = "bb8.domain.tx"
)

// Txer define all methods that can be used for our worker...
type Txer interface {
	InfoByAddress(ctx context.Context, address string) (ResultInfoForADA, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
