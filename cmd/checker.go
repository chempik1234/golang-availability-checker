package main

import (
	"context"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/reporters"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/service"
	"go.uber.org/zap"
	"os"
	"sync"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	configPath := ""
	if len(os.Args) == 2 {
		configPath = os.Args[1]
	} else if len(os.Args) < 2 {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "not enough arguments")
		return
	}
	if configPath == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "please provide a file path argument to use as config")
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "trying to parse config")
	parsedService, err := parsed.ParseService(configPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "trying to set up checkers")
	checkersToRun, err := service.CreateCheckers(ctx, parsedService)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "creating reporters", zap.Int("count", len(checkersToRun)))
	reportersList := make([]reporters.Reporter, len(checkersToRun))

	defer func() {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "closing reporters")
		for _, reporter := range reportersList {
			reporter.Close()
		}
	}()

	wg := &sync.WaitGroup{}

	for index, checker := range checkersToRun {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Creating reporter", zap.String("name", checker.Name), zap.Int("index", index))

		wg.Add(1)
		reportersList[index] = *(reporters.NewReporter(ctx, checker))

		go reportersList[index].Run(wg)
	}

	wg.Wait()
}
