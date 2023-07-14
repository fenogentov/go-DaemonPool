# go-DaemonPool

В пакете реализован запуск демонов с пулом воркеров.
В бесконечном цикле запускается демон, который запускает заданное количество воркеров.
Работа воркеров имеет функционал грецефул шудаун и восстановления при панике.
Определна ошибка ErrNumberWorkers

Пакет разработан для запуска пула демонов
`
jobA := func(ctx context.Contex){

}
func RunDaemons(ctx context.Context, s *Service) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return daemons.WorkersDo(gCtx, "rwAttributes", 5, s.rwAttributes) })
	g.Go(func() error { return daemons.WorkersDo(gCtx, "rwUnitStatus", 1, s.rwUnitStatus) })
	g.Go(func() error { return daemons.WorkersDo(gCtx, "rwConfigObject", 1, s.rwConfigObject) })

	err := g.Wait()
	if gCtx.Err() != nil || err != nil {
		slog.Error(err.Error())
		slog.Error(gCtx.Err().Error())
		slog.Info("shutdown all daemons")
		return
	}
}
`