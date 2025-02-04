package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"

	"github.com/btschwartz12/clips/server"
)

type arguments struct {
	Port       int    `short:"p" long:"port" description:"Port to listen on" default:"8000"`
	VarDir     string `long:"var-dir" env:"CLIPS_VAR_DIR" description:"Directory for storing temporary files" default:"var"`
	ConfigFile string `long:"config-file" env:"CLIPS_CONFIG_FILE" description:"Path to config file" required:"true"`
}

var args arguments

func main() {
	_, err := flags.Parse(&args)
	if err != nil {
		panic(fmt.Errorf("error parsing flags: %s", err))
	}

	var l *zap.Logger
	l, _ = zap.NewDevelopment()
	logger := l.Sugar()

	s := &server.Server{}
	err = s.Init(logger, args.VarDir, args.ConfigFile)
	if err != nil {
		logger.Fatalw("error initializing server", "error", err)
	}
	defer s.Teardown()

	r := chi.NewRouter()
	r.Mount("/", s.Router())

	errChan := make(chan error, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infow("starting http server", "port", args.Port)
		errChan <- http.ListenAndServe(fmt.Sprintf(":%d", args.Port), r)
	}()

	select {
	case err = <-errChan:
		logger.Fatalw("http server failed", "error", err)
	case <-quit:
		logger.Infow("shutting down gracefully")
	}
}
