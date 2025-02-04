package media

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getUpdateTime(t *testing.T, s string) UpdateTime {
	tm, err := time.Parse("15:04", s)
	assert.NoError(t, err)
	ut := UpdateTime{Time: tm}
	return ut
}

func TestNewConfig(t *testing.T) {
	yamlBytes, err := os.ReadFile("example.yaml")
	assert.NoError(t, err)

	config, err := NewConfig(yamlBytes)
	assert.NoError(t, err)

	assert.Len(t, config.Medias, 3)
	media1 := config.Medias[0]
	assert.Equal(t, "Media 1", media1.Name)
	assert.Equal(t, "sample_media/media1", media1.Path)
	assert.True(t, media1.Guessable)
	assert.Equal(t, time.Duration(10000000000), media1.ClipDuration.Duration)
	assert.Equal(t, []UpdateTime{getUpdateTime(t, "10:00"), getUpdateTime(t, "12:00")}, media1.UpdateTimes)
	assert.Equal(t, []string{"sample_media/media1/s1/101.mp4", "sample_media/media1/s1/102.mp4", "sample_media/media1/s2/201.mp4"}, media1.availableVideos)
	assert.True(t, media1.CurrentClip.IsAbsent())

	media2 := config.Medias[1]
	assert.Equal(t, "Media 2", media2.Name)
	assert.Equal(t, "sample_media/media2", media2.Path)
	assert.False(t, media2.Guessable)
	assert.Equal(t, time.Duration(60000000000), media2.ClipDuration.Duration)
	assert.Equal(t, []UpdateTime{getUpdateTime(t, "20:00")}, media2.UpdateTimes)
	assert.Equal(t, []string{"sample_media/media2/s1/101.mp4", "sample_media/media2/s1/102.mp4", "sample_media/media2/s1/103.mp4"}, media2.availableVideos)
	assert.True(t, media2.CurrentClip.IsAbsent())

	media3 := config.Medias[2]
	assert.Equal(t, "Media 3", media3.Name)
	assert.Equal(t, "sample_media/media3", media3.Path)
	assert.False(t, media3.Guessable)
	assert.Equal(t, time.Duration(30000000000), media3.ClipDuration.Duration)
	assert.Equal(t, []UpdateTime{getUpdateTime(t, "15:00")}, media3.UpdateTimes)
	assert.Equal(t, []string{"sample_media/media3/vid1.mp4", "sample_media/media3/vid2.mp4"}, media3.availableVideos)
}
