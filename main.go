package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"

	"github.com/btschwartz12/clips/server"
)

type arguments struct {
	Port      int    `short:"p" long:"port" description:"Port to listen on" default:"8000"`
	MediaDir  string `long:"media-dir" env:"CLIPS_MEDIA_DIR" description:"Directory for media files" required:"true"`
	VarDir    string `long:"var-dir" env:"CLIPS_VAR_DIR" description:"Directory for storing temporary files" default:"var"`
	Duration  string `long:"duration" env:"CLIPS_DURATION" description:"Duration of clips" default:"30s"`
	TimeOfDay string `long:"time-of-day" env:"CLIPS_TIME_OF_DAY" description:"Time of day to update clip"`
}

var args arguments

func main() {
	_, err := flags.Parse(&args)
	if err != nil {
		panic(fmt.Errorf("error parsing flags: %s", err))
	}
	if args.MediaDir == "" {
		panic("media dir is required")
	}
	duration, err := time.ParseDuration(args.Duration)
	if err != nil {
		panic(fmt.Errorf("error parsing duration: %s", err))
	}
	tod, err := time.Parse("15:04:05", args.TimeOfDay)
	if err != nil {
		panic(fmt.Errorf("error parsing time of day: %s", err))
	}

	var l *zap.Logger
	l, _ = zap.NewDevelopment()
	logger := l.Sugar()

	s := &server.Server{}
	err = s.Init(logger, args.MediaDir, args.VarDir, duration, tod)
	if err != nil {
		logger.Fatalw("Error initializing server", "error", err)
	}
	defer s.Teardown()

	r := chi.NewRouter()
	r.Mount("/", s.Router())

	errChan := make(chan error)
	go func() {
		logger.Infow("starting http server", "port", args.Port)
		errChan <- http.ListenAndServe(fmt.Sprintf(":%d", args.Port), r)
	}()
	err = <-errChan
	logger.Fatalw("http server failed", "error", err)
}
