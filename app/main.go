package main

import (
	"github.com/halilylm/microservice/server"
	"go.elastic.co/ecszap"
	_ "go.uber.org/automaxprocs" // for docker container
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var release string

func main() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer func() {
		_ = logger.Sync()
	}()
	logger = logger.With(zap.String("release", release))
	srv := server.New(&server.Options{
		Host:   "0.0.0.0",
		Port:   8080,
		Logger: logger,
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
