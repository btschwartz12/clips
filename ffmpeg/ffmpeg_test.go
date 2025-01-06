package ffmpeg

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sampleVideo = "sample_videos/in1.mp4"

func TestGetVideoDuration(t *testing.T) {
	expectedDuration := time.Duration(7.035646 * float64(time.Second))
	duration, err := getVideoDuration(sampleVideo)
	assert.NoError(t, err)
	assert.Equal(t, expectedDuration, duration)
}

func TestGenerateClip(t *testing.T) {
	outputPath := "sample_videos/out.mp4"
	clipDuration := time.Duration(3 * float64(time.Second))
	err := GenerateClip(sampleVideo, outputPath, clipDuration)
	assert.NoError(t, err)
	duration, err := getVideoDuration(outputPath)
	assert.NoError(t, err)
	const tolerance = 100 * time.Millisecond
	assert.InDelta(t, clipDuration.Seconds(), duration.Seconds(), tolerance.Seconds(), "Durations differ beyond acceptable tolerance")
	err = os.Remove(outputPath)
	assert.NoError(t, err)
}
