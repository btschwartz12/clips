package media

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/samber/mo"
	"gopkg.in/yaml.v3"
)

var (
	supportedFormats = []string{".mp4"}
)

type Clip struct {
	Name string
	Path string
}

type Duration struct {
	time.Duration
}
type UpdateTime struct {
	time.Time
}

type Media struct {
	Name            string       `yaml:"name"`
	Path            string       `yaml:"path"`
	ClipDuration    Duration     `yaml:"clip-duration"`
	UpdateTimes     []UpdateTime `yaml:"update-times"`
	Guessable       bool         `yaml:"guessable"`
	availableVideos []string
	job             mo.Option[gocron.Job]
	CurrentClip     mo.Option[*Clip]
}

type Config struct {
	Medias []Media `yaml:"medias"`
}

func NewConfig(yamlBytes []byte) (*Config, error) {
	c := &Config{}
	if err := yaml.Unmarshal(yamlBytes, c); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	if len(c.Medias) == 0 {
		return nil, fmt.Errorf("no media files found in config file")
	}

	for i := range c.Medias {
		err := c.Medias[i].validate()
		if err != nil {
			return nil, fmt.Errorf("error validating media %d: %w", i, err)
		}
	}
	return c, nil
}

func (m *Media) validate() error {
	if m.Name == "" {
		return fmt.Errorf("media name is required")
	}
	if _, err := os.Stat(m.Path); os.IsNotExist(err) {
		return fmt.Errorf("media path does not exist: %s", m.Path)
	}
	if len(m.UpdateTimes) == 0 {
		return fmt.Errorf("at least one update time is required")
	}
	m.CurrentClip = mo.None[*Clip]()
	err := filepath.Walk(m.Path, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error within walk: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		for _, f := range supportedFormats {
			if strings.HasSuffix(strings.ToLower(info.Name()), f) {
				m.availableVideos = append(m.availableVideos, path)
				break
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}
	return nil
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("error decoding duration: %w", err)
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

func (u *UpdateTime) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("error decoding update time: %w", err)
	}
	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return fmt.Errorf("error parsing update time: %w", err)
	}
	u.Time = parsed
	return nil
}

func (u UpdateTime) MarshalYAML() (interface{}, error) {
	return u.Format("15:04"), nil
}
