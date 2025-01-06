package ffmpeg

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"os/exec"
)

func getVideoDuration(episodePath string) (time.Duration, error) {
	if _, err := os.Stat(episodePath); os.IsNotExist(err) {
		return 0, fmt.Errorf("file not found: %s", episodePath)
	}
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		episodePath,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("error running ffprobe: %w", err)
	}

	durationString := strings.TrimSpace(stdout.String())
	duration, err := time.ParseDuration(durationString + "s")
	if err != nil {
		return 0, fmt.Errorf("error parsing duration: %w", err)
	}

	return duration, nil
}

func GenerateClip(inputPath string, outputPath string, clipDuration time.Duration) error {
	videoDuration, err := getVideoDuration(inputPath)
	if err != nil {
		return fmt.Errorf("error getting duration: %w", err)
	}
	if clipDuration > videoDuration {
		return fmt.Errorf("clip duration is longer than video duration")
	}

	rangeEnd := videoDuration - clipDuration
	startTime := time.Duration(rand.Int63n(int64(rangeEnd)))

	cmd := exec.Command("ffmpeg",
		"-ss", fmt.Sprintf("%02d:%02d:%02d", int(startTime.Hours()), int(startTime.Minutes())%60, int(startTime.Seconds())%60),
		"-i", inputPath,
		"-t", fmt.Sprintf("%02d:%02d:%02d", int(clipDuration.Hours()), int(clipDuration.Minutes())%60, int(clipDuration.Seconds())%60),
		"-c", "copy",
		outputPath,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running ffmpeg: %w. stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("output file not found: %s", outputPath)
	}

	return nil
}
