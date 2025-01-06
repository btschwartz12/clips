package server

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var est *time.Location

func init() {
	var err error
	est, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Errorf("error loading location: %s", err))
	}
}

func (s *Server) addJob(jobFunc any, timeOfDay time.Time) (gocron.Job, error) {
	hours, minutes, seconds, err := parseToDuration(timeOfDay)
	if err != nil {
		return nil, fmt.Errorf("error parsing time to duration: %w", err)
	}
	j, err := s.scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(hours, minutes, seconds),
			),
		),
		gocron.NewTask(
			jobFunc,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating job: %w", err)
	}
	return j, nil
}

func parseToDuration(t time.Time) (hours, minutes, seconds uint, err error) {
	h := t.Hour()
	m := t.Minute()
	s := t.Second()

	if m < 0 || m >= 60 || s < 0 || s >= 60 {
		return 0, 0, 0, fmt.Errorf("invalid time values: hours=%d, minutes=%d, seconds=%d", h, m, s)
	}

	return uint(h), uint(m), uint(s), nil
}

func (s *Server) startCron(timeOfDay time.Time) error {
	sch, err := gocron.NewScheduler(gocron.WithLocation(est))
	if err != nil {
		return fmt.Errorf("error creating scheduler: %w", err)
	}
	s.scheduler = sch

	s.job, err = s.addJob(s.updateClipAsync, timeOfDay)
	if err != nil {
		return fmt.Errorf("error creating job: %w", err)
	}
	s.scheduler.Start()

	s.logger.Infow("cron started, running job now")
	err = s.job.RunNow()
	if err != nil {
		return fmt.Errorf("error running job now: %w", err)
	}
	return nil
}
