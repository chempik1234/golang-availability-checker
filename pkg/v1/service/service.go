package service

import (
	"context"
	"fmt"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/checkers"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
)

func CreateCheckers(ctx context.Context, service *parsed.Service) ([]checkers.ServiceChecker, error) {
	result := make([]checkers.ServiceChecker, len(service.Services))
	for index, checker := range service.Services {
		var repo checkers.CheckerRepository
		switch checker.Protocol {
		case "http":
			repo, _ = checkers.NewCheckerRepositoryHTTP(checker)
		default:
			return nil, fmt.Errorf("expected http, got %s", checker.Protocol)
		}
		result[index] = *checkers.NewServiceChecker(
			ctx, checker.Host, checker.Protocol, checker.Port, checker.IntervalSeconds, checker.TimeoutSeconds, repo,
		)
	}
	return result, nil
}
