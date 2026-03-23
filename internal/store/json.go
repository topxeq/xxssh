package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/topxeq/xxssh/internal/config"
)

type Store struct {
	filePath string
}

func NewStore() (*Store, error) {
	var dir string

	if runtime.GOOS == "windows" {
		// Windows: %USERPROFILE%\.xxssh
		dir = os.Getenv("USERPROFILE")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			dir = home
		}
	} else {
		// Unix: ~/.xxssh
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = home
	}

	dir = filepath.Join(dir, ".xxssh")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Store{
		filePath: filepath.Join(dir, "servers.json"),
	}, nil
}

func (s *Store) Load() (*config.StoresConfig, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &config.StoresConfig{}, nil
		}
		return nil, err
	}
	var cfg config.StoresConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Store) Save(cfg *config.StoresConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}
