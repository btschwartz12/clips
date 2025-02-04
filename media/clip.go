package media

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/mo"

	"github.com/btschwartz12/clips/ffmpeg"
)

func (m *Media) updateClip(outDir string) error {
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(m.availableVideos))))
	if err != nil {
		return fmt.Errorf("error choosing video: %w", err)
	}
	chosenFile := m.availableVideos[idx.Int64()]
	outputPath := filepath.Join(outDir, fmt.Sprintf("%s.mp4", uuid.New()))

	err = ffmpeg.GenerateClip(chosenFile, outputPath, m.ClipDuration.Duration)
	if err != nil {
		return fmt.Errorf("error generating clip: %w", err)
	}

	if m.CurrentClip.IsPresent() {
		err = os.Remove(m.CurrentClip.MustGet().Path)
		if err != nil {
			return fmt.Errorf("error removing old clip: %w", err)
		}
	}
	filename := filepath.Base(chosenFile)
	newClip := Clip{
		Name: strings.TrimSuffix(filename, filepath.Ext(filename)),
		Path: outputPath,
	}
	m.CurrentClip = mo.Some[*Clip](&newClip)
	return nil
}
