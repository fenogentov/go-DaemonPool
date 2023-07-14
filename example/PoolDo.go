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

var jobA = func(ctx context.Context) {
	fmt.Println("jobA")
	time.Sleep(10 * time.Second)
}

func jobB(ctx context.Context) {
	fmt.Println("jobB")
	time.Sleep(12 * time.Second)
}

type job struct {
	name string
}

func (j *job) execute(ctx context.Context) {
	fmt.Println("job.execute")
	time.Sleep(14 * time.Second)
}

func jobMain(ctx context.Context) {
	<-ctx.Done()
	fmt.Println("shutdown jobMain")
}

func RunDaemons(ctx context.Context) {
	jobC := job{name: "jobC"}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return daemons.PoolDo(gCtx, "jobA", 2, jobA) })
	g.Go(func() error { return daemons.PoolDo(gCtx, "jobB", 2, jobB) })
	g.Go(func() error { return daemons.PoolDo(gCtx, "jobC", 2, jobC.execute) })

	if err := g.Wait(); err != nil {
		slog.Error(err.Error())
	}
	if gCtx.Err() != nil {
		slog.Error(gCtx.Err().Error())
	}

	slog.Info("shutdown all daemons")
}

func Example() {
	ctx, stop := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT)
	defer stop()

	go jobMain(ctx)

	RunDaemons(ctx)

	// Ounput:
	// 	job.execute
	// jobB
	// jobA
	// job.execute
	// jobA
	// jobB
	// jobA
	// jobA
	// ^C
	// shutdown jobMain
	// 2023/07/14 12:09:30 INFO shutdown jobC
	// 2023/07/14 12:09:30 INFO shutdown jobB
	// 2023/07/14 12:09:30 INFO shutdown jobA
	// 2023/07/14 12:09:30 ERROR jobC: context canceled
	// 2023/07/14 12:09:30 ERROR context canceled
	// 2023/07/14 12:09:30 INFO shutdown all daemons
}
