package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Server struct {
	router         *chi.Mux
	logger         *zap.SugaredLogger
	mediaDir       string
	varDir         string
	availableFiles []string
	currentClip    *Clip
	clipDuration   time.Duration
	scheduler      gocron.Scheduler
	job            gocron.Job
}

func (s *Server) Init(
	logger *zap.SugaredLogger,
	mediaDir string,
	varDir string,
	clipDuration time.Duration,
	timeOfDay time.Time,
) error {
	s.logger = logger
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		return fmt.Errorf("media directory does not exist: %s", mediaDir)
	}
	s.mediaDir = mediaDir

	if _, err := os.Stat(varDir); os.IsNotExist(err) {
		err = os.Mkdir(varDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating var directory: %w", err)
		}
	}
	s.varDir = varDir
	s.clipDuration = clipDuration
	s.router = chi.NewRouter()

	s.router.Get("/", s.home)
	s.router.Get("/clip", s.serveClip)
	s.router.Get("/health", s.health)

	files, err := s.getVideoFiles()
	if err != nil {
		return fmt.Errorf("error getting video files: %w", err)
	}
	s.availableFiles = files
	logger.Infow("available video files", "numFiles", len(files))

	err = s.startCron(timeOfDay)
	if err != nil {
		return fmt.Errorf("error starting cron: %w", err)
	}
	nextRun, err := s.job.NextRun()
	if err != nil {
		return fmt.Errorf("error getting next run time: %w", err)
	}
	logger.Infow("next run time", "time", nextRun)

	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) getVideoFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(s.mediaDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		if err != nil {
			s.logger.Errorw("error walking directory", "error", err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".mp4") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return files, nil
}

func (s *Server) Teardown() {
	err := s.scheduler.Shutdown()
	if err != nil {
		s.logger.Errorw("error shutting down scheduler", "error", err)
	}
}
