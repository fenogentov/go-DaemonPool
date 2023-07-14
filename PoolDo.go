package godaemonpool

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

var ErrNumberWorkers = errors.New("invalid number of workers")

// PoolDo
func PoolDo(ctx context.Context, name string, maxWorkers int, job func(ctx context.Context)) error {
	if maxWorkers < 1 {
		return ErrNumberWorkers
	}

	ch := make(chan struct{}, maxWorkers)

	for {
		select {
		case ch <- struct{}{}:
			go func(ctx context.Context, ch chan struct{}) {
				defer func() {
					if p := recover(); p != nil {
						msg := fmt.Sprintf("panic '%s'", name)
						slog.Error(msg, "recover", p, "stack", string(debug.Stack()))
					}
					<-ch
				}()
				job(ctx)
			}(ctx, ch)

		case <-ctx.Done():
			slog.Info(fmt.Sprintf("shutdown %s", name))
			return fmt.Errorf("%s: %w", name, ctx.Err())

		default:
			continue
		}
	}
}
