package server

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btschwartz12/clips/ffmpeg"
	"github.com/google/uuid"
)

type Clip struct {
	Name string
	Path string
}



func (s *Server) updateClip() error {

	chosenFile := s.availableFiles[rand.Int63n(int64(len(s.availableFiles)))]
	fileName := filepath.Base(chosenFile)
	outputPath := filepath.Join(s.varDir, fmt.Sprintf("%s.mp4", uuid.New()))

	s.logger.Infow("generating clip", "file", chosenFile, "outputPath", outputPath)
	err := ffmpeg.GenerateClip(chosenFile, outputPath, s.clipDuration)
	if err != nil {
		return fmt.Errorf("error generating clip: %w", err)
	}
	s.logger.Infow("clip generated", "file", chosenFile, "outputPath", outputPath)
	
	err = filepath.Walk(s.varDir, func(path string, info os.FileInfo, err error) error {
		if path == outputPath {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		return os.Remove(path)
	})
	if err != nil {
		return fmt.Errorf("error removing old clips: %w", err)
	}
	s.currentClip = &Clip{
		Name: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		Path: outputPath,
	}
	return nil
}

func (s *Server) updateClipDaily() {
	err := s.updateClip()
	if err != nil {
		s.logger.Errorw("error updating clip", "error", err)
	}

	now := time.Now().UTC()
	year, month, day := now.Date()
	nextUpdate := time.Date(year, month, day, 17, 0, 0, 0, time.UTC)
	if now.After(nextUpdate) {
		nextUpdate = nextUpdate.Add(24 * time.Hour)
	}
	timeUntilUpdate := nextUpdate.Sub(now)
	s.logger.Infow("waiting until next clip update", "timeUntilUpdate", timeUntilUpdate)
	time.Sleep(timeUntilUpdate)

	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			err := s.updateClip()
			if err != nil {
				s.logger.Errorw("error updating clip", "error", err)
				continue
			}
			s.logger.Infow("clip updated")
		}
	}
}