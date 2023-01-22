package main

import (
	"github.com/halilylm/microservice/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var release string

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	logger = logger.With(zap.String("release", release))
	sugar := logger.Sugar()
	srv := server.New(&server.Options{
		Host:   "localhost",
		Port:   8080,
		Logger: sugar,
	})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, os.Interrupt)
	go func() {
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()
	<-quit
	if err := srv.Stop(); err != nil {
		panic(err)
	}
}
