package media

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/samber/mo"
	"go.uber.org/zap"
)

type durationTuple struct {
	hours   uint
	minutes uint
	seconds uint
}

func (m *Media) getCronTimes() (gocron.AtTimes, error) {
	if len(m.UpdateTimes) == 0 {
		return nil, fmt.Errorf("no update times found")
	}
	tuples := make([]durationTuple, len(m.UpdateTimes))
	for i, ut := range m.UpdateTimes {
		h := ut.Hour()
		m := ut.Minute()
		s := ut.Second()

		if m < 0 || m >= 60 || s < 0 || s >= 60 {
			return nil, fmt.Errorf("invalid time values: hours=%d, minutes=%d, seconds=%d", h, m, s)
		}
		tuples[i] = durationTuple{uint(h), uint(m), uint(s)}
	}
	atTimes := make([]gocron.AtTime, len(tuples))
	for i, t := range tuples {
		atTimes[i] = gocron.NewAtTime(t.hours, t.minutes, t.seconds)
	}
	return gocron.NewAtTimes(atTimes[0], atTimes[1:]...), nil
}

func (m *Media) GetJob(
	logger *zap.SugaredLogger,
	outDir string,
) (gocron.JobDefinition, gocron.Task, error) {
	atTimes, err := m.getCronTimes()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting cron times: %w", err)
	}
	jd := gocron.DailyJob(
		1,
		atTimes,
	)
	task := gocron.NewTask(
		m.UpdateClipAsync,
		logger,
		outDir,
	)
	return jd, task, nil
}

func (m *Media) UpdateClipAsync(logger *zap.SugaredLogger, outDir string) {
	err := m.updateClip(outDir)
	if err != nil {
		logger.Errorw("error updating clip", "error", err)
	}
}

func (m *Media) SetJob(job gocron.Job) {
	m.job = mo.Some(job)
}

func (m *Media) RunJobNow() error {
	if m.job.IsAbsent() {
		return fmt.Errorf("job is not present")
	}
	return m.job.MustGet().RunNow()
}

func (m *Media) GetNextRun() (time.Time, error) {
	if m.job.IsAbsent() {
		return time.Time{}, fmt.Errorf("job is not present")
	}
	return m.job.MustGet().NextRun()
}
