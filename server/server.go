package server

import (
	"fmt"
	"os"

	"github.com/btschwartz12/clips/media"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Server struct {
	router    *chi.Mux
	logger    *zap.SugaredLogger
	config    *media.Config
	varDir    string
	scheduler gocron.Scheduler
}

func (s *Server) Init(
	logger *zap.SugaredLogger,
	varDir string,
	configPath string,
) error {
	s.logger = logger

	yamlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	s.config, err = media.NewConfig(yamlBytes)
	if err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	if _, err := os.Stat(varDir); os.IsNotExist(err) {
		err = os.Mkdir(varDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating var directory: %w", err)
		}
	}
	s.varDir = varDir

	s.router = chi.NewRouter()
	s.router.Get("/", s.home)
	s.router.Get("/clip/{name}", s.serveClip)
	s.router.Get("/health", s.health)

	err = s.startCron()
	if err != nil {
		return fmt.Errorf("error starting cron: %w", err)
	}
	return nil
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) Teardown() {
	for _, m := range s.config.Medias {
		if m.CurrentClip.IsPresent() {
			err := os.Remove(m.CurrentClip.MustGet().Path)
			if err != nil {
				s.logger.Errorw("error deleting media file", "error", err)
			}
		}
	}
	err := s.scheduler.Shutdown()
	if err != nil {
		s.logger.Errorw("error shutting down scheduler", "error", err)
	}
}
