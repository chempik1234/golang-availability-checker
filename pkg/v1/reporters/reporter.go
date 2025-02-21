package reporters

import (
	"context"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/checkers"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/types"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/utils"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/webhooks"
	"go.uber.org/zap"
	"sync"
)

type Reporter struct {
	ctx             context.Context
	checker         checkers.ServiceChecker
	resultChannel   types.ResultChan
	webhooksManager webhooks.Manager
	Name            string
}

func NewReporter(ctx context.Context, checker checkers.ServiceChecker, webhooksManager webhooks.Manager) *Reporter {
	return &Reporter{
		ctx:             ctx,
		checker:         checker,
		resultChannel:   types.NewResultChan(),
		webhooksManager: webhooksManager,
		Name:            checker.Name,
	}
}

func (r *Reporter) Close() {
	logger.GetLoggerFromCtx(r.ctx).Info(r.ctx, "closing reporter", zap.String("reporter", r.Name))
	close(r.resultChannel)
}

func (r *Reporter) Run(wg *sync.WaitGroup) {
	defer func() {
		_ = recover()
		logger.GetLoggerFromCtx(r.ctx).Info(r.ctx, "reporter goroutine exited", zap.String("reporter", r.Name))
		close(r.resultChannel)
	}()
	defer wg.Done()

	logger.GetLoggerFromCtx(r.ctx).Info(
		r.ctx, "listening begins",
		zap.String("URL", utils.Url(r.checker.Protocol, r.checker.Host, r.checker.Port)),
	)
	go r.checker.Run(r.resultChannel)
	for result := range r.resultChannel {
		logger.GetLoggerFromCtx(r.ctx).Info(
			r.ctx,
			"tried requesting a service",
			zap.String("name", r.Name),
			zap.Bool("result", result),
		)
		httpBody := webhooks.FormHTTPBody(r.Name, result)
		r.webhooksManager.Report(httpBody)
	}
}
