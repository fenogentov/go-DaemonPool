package godaemonpool

import (
	"context"
	"fmt"
	"os/signal"
	"time"

	daemons "github.com/fenogentov/go-DaemonPool"

	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

var jobAE = func(ctx context.Context) error {
	slog.Info("jobAE")
	time.Sleep(10 * time.Second)
	return nil
}

func jobBE(ctx context.Context) error {
	slog.Info("jobBE")
	time.Sleep(12 * time.Second)
	return fmt.Errorf("failed jobBE")
}

type jobE struct {
	name string
}

func (j *jobE) execute(ctx context.Context) error {
	slog.Info("jobE.execute")
	time.Sleep(14 * time.Second)
	return nil
}

func jobEMain(ctx context.Context) {
	<-ctx.Done()
	slog.Info("shutdown jobEMain")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	defer stop()

	go jobEMain(ctx)

	RunDaemonsE(ctx)
}

func RunDaemonsE(ctx context.Context) {
	jobCE := jobE{name: "jobCE"}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return daemons.PoolWithErrorDo(gCtx, "jobEA", 2, jobAE) })
	g.Go(func() error { return daemons.PoolWithErrorDo(gCtx, "jobEB", 2, jobBE) })
	g.Go(func() error { return daemons.PoolWithErrorDo(gCtx, "jobEC", 2, jobCE.execute) })

	if err := g.Wait(); err != nil {
		slog.Error(err.Error())
	}
	if gCtx.Err() != nil {
		slog.Error(gCtx.Err().Error())
	}

	slog.Info("shutdown all daemons")

	// Output:
	// 2023/07/14 12:23:00 INFO job.execute
	// 2023/07/14 12:23:00 INFO jobA
	// 2023/07/14 12:23:00 INFO jobB
	// 2023/07/14 12:23:00 INFO jobB
	// 2023/07/14 12:23:00 INFO jobA
	// 2023/07/14 12:23:00 INFO job.execute
	// 2023/07/14 12:23:10 INFO jobA
	// 2023/07/14 12:23:10 INFO jobA
	// 2023/07/14 12:23:12 INFO shutdown jobEA
	// 2023/07/14 12:23:12 INFO shutdown jobEC
	// 2023/07/14 12:23:12 ERROR jobEB: failed jobBE
	// 2023/07/14 12:23:12 ERROR context canceled
	// 2023/07/14 12:23:12 INFO shutdown all daemons
	// 2023/07/14 12:23:12 INFO shutdown jobMain
}
