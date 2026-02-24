package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	appName    = "tui-timer"
	configFile = "config.yaml"
)

type SoundsConfig struct {
	Tick   bool `yaml:"tick"`
	Finish bool `yaml:"finish"`
	Break  bool `yaml:"break"`
}

type VoiceMessages struct {
	WorkDone  string `yaml:"work_done"`
	BreakDone string `yaml:"break_done"`
	Start     string `yaml:"start"`
}

type VoiceConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Voice    string        `yaml:"voice"`
	Messages VoiceMessages `yaml:"messages"`
}

type Config struct {
	WorkDuration     time.Duration `yaml:"-"`
	ShortBreak       time.Duration `yaml:"-"`
	LongBreak        time.Duration `yaml:"-"`
	CyclesBeforeLong int           `yaml:"cycles_before_long"`

	// YAML string fields for serialization
	WorkDurationStr string `yaml:"work_duration"`
	ShortBreakStr   string `yaml:"short_break"`
	LongBreakStr    string `yaml:"long_break"`

	Sounds SoundsConfig `yaml:"sounds"`
	Voice  VoiceConfig  `yaml:"voice"`
}

func DefaultConfig() *Config {
	return &Config{
		WorkDuration:     25 * time.Minute,
		ShortBreak:       5 * time.Minute,
		LongBreak:        15 * time.Minute,
		CyclesBeforeLong: 4,
		WorkDurationStr:  "25m",
		ShortBreakStr:    "5m",
		LongBreakStr:     "15m",
		Sounds: SoundsConfig{
			Tick:   true,
			Finish: true,
			Break:  true,
		},
		Voice: VoiceConfig{
			Enabled: true,
			Voice:   "Samantha",
			Messages: VoiceMessages{
				WorkDone:  "Work session finished",
				BreakDone: "Break finished",
				Start:     "Focus time started",
			},
		},
	}
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", appName), nil
}

func ConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := Save(cfg); err != nil {
				return cfg, fmt.Errorf("creating default config: %w", err)
			}
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.parseDurations(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, configFile)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (c *Config) parseDurations() error {
	var err error
	if c.WorkDurationStr != "" {
		c.WorkDuration, err = time.ParseDuration(c.WorkDurationStr)
		if err != nil {
			return fmt.Errorf("invalid work_duration: %w", err)
		}
	}
	if c.ShortBreakStr != "" {
		c.ShortBreak, err = time.ParseDuration(c.ShortBreakStr)
		if err != nil {
			return fmt.Errorf("invalid short_break: %w", err)
		}
	}
	if c.LongBreakStr != "" {
		c.LongBreak, err = time.ParseDuration(c.LongBreakStr)
		if err != nil {
			return fmt.Errorf("invalid long_break: %w", err)
		}
	}
	return nil
}

// ApplyCLIFlags parses CLI flags and overrides config values.
func (c *Config) ApplyCLIFlags(args []string) error {
	fs := flag.NewFlagSet("tui-timer", flag.ContinueOnError)

	work := fs.String("work", "", "work duration (e.g. 50m)")
	shortBreak := fs.String("short-break", "", "short break duration")
	longBreak := fs.String("long-break", "", "long break duration")
	voice := fs.String("voice", "", "macOS voice name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *work != "" {
		d, err := time.ParseDuration(*work)
		if err != nil {
			return fmt.Errorf("invalid --work: %w", err)
		}
		c.WorkDuration = d
		c.WorkDurationStr = *work
	}
	if *shortBreak != "" {
		d, err := time.ParseDuration(*shortBreak)
		if err != nil {
			return fmt.Errorf("invalid --short-break: %w", err)
		}
		c.ShortBreak = d
		c.ShortBreakStr = *shortBreak
	}
	if *longBreak != "" {
		d, err := time.ParseDuration(*longBreak)
		if err != nil {
			return fmt.Errorf("invalid --long-break: %w", err)
		}
		c.LongBreak = d
		c.LongBreakStr = *longBreak
	}
	if *voice != "" {
		c.Voice.Voice = *voice
	}

	return nil
}
