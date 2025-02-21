package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/golang-availability-checker/pkg/logger"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/checkersmanager"
	"github.com/chempik1234/golang-availability-checker/pkg/v1/parsed"
	"os"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	configPath := ""
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	} else if len(os.Args) < 2 {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "not enough arguments")
	}
	if configPath == "" {
		fmt.Print("please provide a file path argument to use as config: ")
		_, err := fmt.Scanln(&configPath)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Error(ctx, "file path wasn't provided")
			return
		}
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "trying to parse config")
	parsedService, err := parsed.ParseService(configPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "trying to set up checkers")
	service, err := checkersmanager.CreateCheckersManager(ctx, parsedService)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, err.Error())
		return
	}

	service.Run()
}
