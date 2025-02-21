package checkers

import (
	"context"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/types"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/utils"
	"net"
	"net/http"
	"time"
)

type CheckerRepositoryHTTP struct {
	transport *http.Transport
	client    *http.Client
	req       *http.Request
}

func NewCheckerRepositoryHTTP(checker parsed.Checker) (CheckerRepositoryHTTP, error) {
	timeout := time.Duration(checker.TimeoutSeconds) * time.Second
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		IdleConnTimeout: timeout,
	}
	url := utils.Url(checker.Protocol, checker.Host, checker.Port)
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return CheckerRepositoryHTTP{}, err
	}
	return CheckerRepositoryHTTP{
		transport: transport,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		req: req,
	}, nil
}

func (c CheckerRepositoryHTTP) Check(ctx context.Context, result types.ResultChan) {
	resp, err := c.client.Do(c.req)
	if err != nil {
		result.WriteFailure()
		return
	}
	if resp.StatusCode != http.StatusOK {
		result.WriteFailure()
		return
	}
	result.WriteSuccess()
}
