package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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
}

func (s *Server) Init(
	logger *zap.SugaredLogger,
	mediaDir string,
	varDir string,
	clipDuration time.Duration,
) error {
	s.logger = logger
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		return fmt.Errorf("media directory does not exist: %s", mediaDir)
	}
	s.mediaDir = mediaDir

	if _, err := os.Stat(varDir); os.IsNotExist(err) {
		return fmt.Errorf("var directory does not exist: %s", varDir)
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

	go s.updateClipDaily()

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
