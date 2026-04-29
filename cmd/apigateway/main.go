package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apigateway/internal/config"
	"apigateway/internal/logger"
	"apigateway/internal/server"

	"go.uber.org/zap"
)

func main() {
	// 1. Load config
	if err := config.Init("config.yaml"); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init logger
	logger.Init(&logger.Config{
		Level:      config.GlobalConfig.Log.Level,
		Filename:   config.GlobalConfig.Log.Filename,
		MaxSize:    config.GlobalConfig.Log.MaxSize,
		MaxBackups: config.GlobalConfig.Log.MaxBackups,
		MaxAge:     config.GlobalConfig.Log.MaxAge,
		Compress:   config.GlobalConfig.Log.Compress,
	})

	// 3. Create and start server
	srv := server.NewServer(&config.GlobalConfig)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Log.Fatal("server error", zap.Error(err))
		}
	}()

	logger.Log.Info("API Gateway started")

	// 4. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("shutdown error", zap.Error(err))
	}

	logger.Log.Info("API Gateway stopped")
}
