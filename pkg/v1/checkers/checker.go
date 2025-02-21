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
	ctx             context.Context
	Name            string
	Host            string
	Protocol        string
	Port            string
	intervalSeconds time.Duration
	timeoutSeconds  time.Duration
	repo            CheckerRepository
}

func NewServiceChecker(
	ctx context.Context,
	name string,
	host string,
	protocol string,
	port string,
	intervalSeconds int,
	timeoutSeconds int,
	repo CheckerRepository) *ServiceChecker {
	return &ServiceChecker{
		ctx:             ctx,
		Name:            name,
		Host:            host,
		Protocol:        protocol,
		Port:            port,
		intervalSeconds: time.Second * time.Duration(intervalSeconds),
		timeoutSeconds:  time.Second * time.Duration(timeoutSeconds),
		repo:            repo,
	}
}

func (c *ServiceChecker) Run(reportChannel types.ResultChan) {
	defer func() {
		logger.GetLoggerFromCtx(c.ctx).Info(c.ctx, "restaring checker", zap.String("name", c.Name))
		_ = recover()
		go c.Run(reportChannel)
	}()

	var requestContext context.Context
	var cancel context.CancelFunc
	var currentResult bool

	ticker := time.NewTicker(c.intervalSeconds)

	for range ticker.C {
		currentResult = false
		requestContext, cancel = context.WithTimeout(context.Background(), c.timeoutSeconds)
		resultChannel := types.NewResultChan()
		go c.repo.Check(requestContext, resultChannel)
	LOOP:
		for {
			select {
			case <-c.ctx.Done():
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
