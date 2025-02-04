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

func (s *Server) addJobs() error {
	for i := range s.config.Medias {
		jd, task, err := s.config.Medias[i].GetJob(s.logger, s.varDir)
		if err != nil {
			return fmt.Errorf("error getting job: %w", err)
		}
		j, err := s.scheduler.NewJob(jd, task)
		if err != nil {
			return fmt.Errorf("error creating job: %w", err)
		}
		s.config.Medias[i].SetJob(j)
	}
	return nil
}

func (s *Server) startCron() error {
	sch, err := gocron.NewScheduler(gocron.WithLocation(est))
	if err != nil {
		return fmt.Errorf("error creating scheduler: %w", err)
	}
	s.scheduler = sch

	err = s.addJobs()
	if err != nil {
		return fmt.Errorf("error adding jobs: %w", err)
	}
	s.scheduler.Start()

	s.logger.Infow("cron started, running jobs now")
	for _, m := range s.config.Medias {
		err := m.RunJobNow()
		if err != nil {
			return fmt.Errorf("error running job now: %w", err)
		}
	}
	s.logger.Infow("jobs have been run initially")
	for _, m := range s.config.Medias {
		nextRun, err := m.GetNextRun()
		if err != nil {
			return fmt.Errorf("error getting next run time: %w", err)
		}
		s.logger.Infow("next run time", "media", m.Name, "time", nextRun)
	}
	return nil
}
