package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds huijata's persisted settings.
type Config struct {
	// GamePath is the Nolla_Games_Noita folder that contains the save slots.
	GamePath string `toml:"game_path"`
}

// SavePath returns the active save slot folder (save00) within GamePath.
func (c *Config) SavePath() string {
	return filepath.Join(c.GamePath, "save00")
}

// Path returns the location of huijata's config file. When HUIJATA_CONFIG_PATH
// is set it is treated as the directory containing config.toml. Otherwise the
// config lives in the user's config dir (on Linux
// ~/.config/huijata/config.toml; on Windows %AppData%\huijata\config.toml).
func Path() (string, error) {
	if p := os.Getenv("HUIJATA_CONFIG_PATH"); p != "" {
		return filepath.Join(p, "config.toml"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "huijata", "config.toml"), nil
}

// SnapshotsDir returns where save snapshots are stored, a snapshots folder
// next to the config file.
func SnapshotsDir() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), "snapshots"), nil
}

// Load reads the config file, returning nil if it does not exist yet.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to disk, creating the config directory if needed.
func (c *Config) Save() error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(c)
}
