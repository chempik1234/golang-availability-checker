package checkersmanager

import (
	"context"
	"fmt"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/checkers"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/reporters"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/webhooks"
	"go.uber.org/zap"
	"sync"
)

type CheckersManager struct {
	ctx       context.Context
	reporters []reporters.Reporter
}

func CreateCheckersManager(ctx context.Context, service *parsed.Service) (*CheckersManager, error) {
	reportersList := make([]reporters.Reporter, len(service.Services))
	for index, checker := range service.Services {
		var repo checkers.CheckerRepository
		switch checker.Protocol {
		case "http":
			repo, _ = checkers.NewCheckerRepositoryHTTP(checker)
		default:
			return &CheckersManager{}, fmt.Errorf("expected http, got %s", checker.Protocol)
		}
		webhooksManager, err := webhooks.NewManager(ctx, checker)
		if err != nil {
			return &CheckersManager{}, err
		}
		reportersList[index] = *reporters.NewReporter(
			ctx,
			*checkers.NewServiceChecker(
				ctx, checker.Name, checker.Host, checker.Protocol, checker.Port, checker.IntervalSeconds, checker.TimeoutSeconds, repo,
			),
			webhooksManager,
		)
	}
	return &CheckersManager{ctx: ctx, reporters: reportersList}, nil
}

func (c *CheckersManager) CheckersCount() int {
	return len(c.reporters)
}

func (c *CheckersManager) Run() {
	logger.GetLoggerFromCtx(c.ctx).Info(c.ctx, "creating reporters", zap.Int("count", c.CheckersCount()))

	defer func() {
		logger.GetLoggerFromCtx(c.ctx).Info(c.ctx, "closing reporters")
		for _, reporter := range c.reporters {
			reporter.Close()
		}
	}()

	wg := &sync.WaitGroup{}

	for index, checker := range c.reporters {
		logger.GetLoggerFromCtx(c.ctx).Info(c.ctx,
			"Creating reporter", zap.String("name", checker.Name), zap.Int("index", index))

		wg.Add(1)

		go c.reporters[index].Run(wg)
	}

	wg.Wait()
}
