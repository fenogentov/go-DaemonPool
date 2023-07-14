# go-DaemonPool

В пакете реализован запуск демонов с пулом воркеров.
В бесконечном цикле запускается демон, который запускает заданное количество воркеров.

Работа воркеров имеет функционал грецефул шудаун и восстановления при панике.

Определна ошибка *ErrNumberWorkers* - если количество воркеров задано меньше 1.

Пакет разработан для запуска пула демонов
``` go
func RunDaemons(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return daemons.PoolDo(gCtx, "jobA", 2, jobA) })
	g.Go(func() error { return daemons.PoolDo(gCtx, "jobB", 2, jobB) })

	if err := g.Wait(); err != nil {
		slog.Error(err.Error())
	}
	if gCtx.Err() != nil {
		slog.Error(gCtx.Err().Error())
	}

	slog.Info("shutdown all daemons")
}
```
Пакет имеет функции:
* `PoolDo` - для работ имеющих сигнатуру *`job(ctx context.Context)`*. Данная фукция подразумевает обработку ошибок внутри воркера
* `PoolWithErrorDo` - для работ имеющих сигнатуру *`job(ctx context.Context)error`*. Данная функция заканчивает работу пула воркеров если в одном из воркеров произошла ошибка.