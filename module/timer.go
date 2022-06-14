package module

import (
	"context"
	"time"
)

type TimerCallback func(ctx context.Context, cancelFunc context.CancelFunc)

func (b *Base) NewTimer(ctxIn context.Context, timeIn time.Duration, fn TimerCallback) context.CancelFunc {
	ctx, cancelFn := context.WithCancel(ctxIn)
	t := time.NewTicker(timeIn)

	go func() {
		for {
			select {
			case <-t.C:
				fn(ctx, cancelFn)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancelFn
}
