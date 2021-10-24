package bb8

import (
	"context"
	"time"
)

type Txer interface {
	InfoByAddress(address string) (ResultInfoForADA, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
