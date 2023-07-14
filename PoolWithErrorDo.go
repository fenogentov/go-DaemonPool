package godaemonpool

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// PoolWithErrorDo запускает в горутинах требуемое количество воркеров
// с функционалом Graceful Shutdown и контролем ошибок в воркерах.
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
