package godaemonpool

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// PoolDo
func PoolWithErrorDo(ctx context.Context, name string, maxWorkers int, job func(ctx context.Context) error) error {
	if maxWorkers < 1 {
		return ErrNumberWorkers
	}

	chWorkers := make(chan struct{}, maxWorkers)
	chError := make(chan error, 1)

	for {
		select {
		case err := <-chError:
			return err

		case <-ctx.Done():
			slog.Info(fmt.Sprintf("shutdown %s", name))
			return fmt.Errorf("%s: %w", name, ctx.Err())

		case chWorkers <- struct{}{}:
			go func(ctx context.Context, chW chan struct{}, chE chan error) {
				defer func() {
					if p := recover(); p != nil {
						msg := fmt.Sprintf("panic '%s'", name)
						slog.Error(msg, "recover", p, "stack", string(debug.Stack()))
					}
					<-chW
				}()

				err := job(ctx)
				if err != nil {
					chE <- fmt.Errorf("%s: %w", name, err)
				}
			}(ctx, chWorkers, chError)

		default:
			continue
		}
	}
}
