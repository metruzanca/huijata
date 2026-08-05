package cmd

import (
	"errors"
	"strings"

	"github.com/metruzanca/huijata/internal/config"
)

// loadConfig reads the config, failing with a hint to run init when missing.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("no config found. Run `huijata init` first.")
	}
	return cfg, nil
}

func validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("description must not be empty")
	}
	if strings.ContainsAny(description, `/\`) {
		return errors.New("description must not contain path separators")
	}
	return nil
}
