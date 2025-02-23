package webhooks

import (
	"bytes"
	"context"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Webhook struct {
	Url        string
	Headers    map[string]string
	StatusCode int
}

type Manager struct {
	ctx       context.Context
	transport *http.Transport
	client    *http.Client
	webhooks  []Webhook
}

func NewManager(ctx context.Context, checker parsed.Checker) (Manager, error) {
	timeout := time.Duration(checker.TimeoutSeconds) * time.Second
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		IdleConnTimeout: timeout,
	}
	webhooks := make([]Webhook, len(checker.Webhooks))
	for i, webhook := range checker.Webhooks {
		webhooks[i] = Webhook{
			Url:        webhook.Url,
			Headers:    webhook.Headers,
			StatusCode: webhook.StatusCode,
		}
	}
	return Manager{
		ctx:       ctx,
		transport: transport,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		webhooks: webhooks,
	}, nil
}

func (w *Manager) Report(httpBody string) {
	for _, webhook := range w.webhooks {
		go func() {
			body := bytes.NewBufferString(httpBody)
			request, err := http.NewRequest(http.MethodPost, webhook.Url, body)
			for k, v := range webhook.Headers {
				request.Header.Set(k, v)
			}
			if err != nil {
				logger.GetLoggerFromCtx(w.ctx).Error(
					w.ctx,
					"couldn't create a webhook request object",
					zap.String("url", webhook.Url),
					zap.Error(err),
				)
			}
			response, err := w.client.Do(request)
			if err != nil {
				logger.GetLoggerFromCtx(w.ctx).Error(
					w.ctx,
					"couldn't send a webhook POST request",
					zap.String("url", webhook.Url),
					zap.Error(err),
				)
				return
			}
			if response.StatusCode != webhook.StatusCode {
				logger.GetLoggerFromCtx(w.ctx).Error(
					w.ctx,
					"webhook doesn't answer 200 on POST request",
					zap.String("url", webhook.Url),
					zap.Error(err),
				)
				return
			}
		}()
	}
}
