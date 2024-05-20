package utils

import (
	"context"
	"time"
)

func NewTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
