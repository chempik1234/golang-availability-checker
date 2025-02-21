package checkers

import (
	"context"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/types"
	"go.uber.org/zap"
	"time"
)

type CheckerRepository interface {
	Check(ctx context.Context, result types.ResultChan)
}

type ServiceChecker struct {
	Name            string
	Ctx             context.Context
	Host            string
	Protocol        string
	Port            string
	IntervalSeconds time.Duration
	TimeoutSeconds  time.Duration
	Repo            CheckerRepository
}

func NewServiceChecker(
	ctx context.Context,
	host string,
	protocol string,
	port string,
	intervalSeconds int,
	timeoutSeconds int,
	repo CheckerRepository) *ServiceChecker {
	return &ServiceChecker{
		Ctx:             ctx,
		Host:            host,
		Protocol:        protocol,
		Port:            port,
		IntervalSeconds: time.Second * time.Duration(intervalSeconds),
		TimeoutSeconds:  time.Second * time.Duration(timeoutSeconds),
		Repo:            repo,
	}
}

func (c *ServiceChecker) Run(reportChannel types.ResultChan) {
	defer func() {
		logger.GetLoggerFromCtx(c.Ctx).Info(c.Ctx, "restaring checker", zap.String("name", c.Name))
		_ = recover()
		go c.Run(reportChannel)
	}()

	var requestContext context.Context
	var cancel context.CancelFunc
	var currentResult bool

	ticker := time.NewTicker(c.IntervalSeconds)

	for range ticker.C {
		currentResult = false
		requestContext, cancel = context.WithTimeout(context.Background(), c.TimeoutSeconds)
		resultChannel := types.NewResultChan()
		go c.Repo.Check(requestContext, resultChannel)
	LOOP:
		for {
			select {
			case <-c.Ctx.Done():
				reportChannel.WriteFailure()
				break LOOP
			case currentResult = <-resultChannel:
				if currentResult {
					reportChannel.WriteSuccess()
				} else {
					reportChannel.WriteFailure()
				}
				break LOOP
			}
		}
		cancel()
	}
}
