package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/topxeq/xxssh/internal/config"
	"github.com/topxeq/xxssh/internal/crypto"
)

type Store struct {
	filePath   string
	masterKey  []byte
	isUnlocked bool
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

// SetMasterPassword sets the master password for encrypting/decrypting passwords
func (s *Store) SetMasterPassword(password string) error {
	cfg, err := s.LoadRaw()
	if err != nil {
		return err
	}

	if cfg.Salt == "" {
		// First time setting password - generate salt and hash
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return err
		}
		cfg.Salt = crypto.Base64Encode(salt)
		cfg.MasterPasswordHash = crypto.HashPassword([]byte(password), salt)
		s.masterKey = crypto.DeriveKey([]byte(password), salt)
		s.isUnlocked = true
	} else {
		// Verify password
		salt := crypto.Base64Decode(cfg.Salt)
		hash := crypto.HashPassword([]byte(password), salt)
		if hash != cfg.MasterPasswordHash {
			return crypto.ErrInvalidPassword
		}
		s.masterKey = crypto.DeriveKey([]byte(password), salt)
		s.isUnlocked = true
	}

	// Re-save with new hash (don't re-encrypt passwords since they use the same key)
	return s.SaveRaw(cfg)
}

// IsLocked returns true if the store requires a master password
func (s *Store) IsLocked() bool {
	cfg, _ := s.LoadRaw()
	return cfg.Salt == ""
}

// Unlock unlocks the store with the given master password
func (s *Store) Unlock(password string) error {
	return s.SetMasterPassword(password)
}

// Lock locks the store (clears the master key)
func (s *Store) Lock() {
	s.masterKey = nil
	s.isUnlocked = false
}

// IsUnlocked returns true if the store is currently unlocked
func (s *Store) IsUnlocked() bool {
	return s.isUnlocked
}

// ChangeMasterPassword changes the master password and re-encrypts all passwords
func (s *Store) ChangeMasterPassword(oldPassword, newPassword string) error {
	cfg, err := s.LoadRaw()
	if err != nil {
		return err
	}

	if cfg.Salt == "" {
		return errors.New("no master password set")
	}

	// Verify old password
	salt := crypto.Base64Decode(cfg.Salt)
	hash := crypto.HashPassword([]byte(oldPassword), salt)
	if hash != cfg.MasterPasswordHash {
		return crypto.ErrInvalidPassword
	}

	// Decrypt all passwords with old key
	oldKey := crypto.DeriveKey([]byte(oldPassword), salt)
	for i := range cfg.Servers {
		if cfg.Servers[i].Password != "" {
			decrypted, err := crypto.Decrypt(cfg.Servers[i].Password, oldKey)
			if err != nil {
				return err
			}
			cfg.Servers[i].Password = decrypted
		}
	}

	// Generate new salt and key
	newSalt, err := crypto.GenerateSalt()
	if err != nil {
		return err
	}
	newKey := crypto.DeriveKey([]byte(newPassword), newSalt)

	// Re-encrypt all passwords with new key
	for i := range cfg.Servers {
		if cfg.Servers[i].Password != "" {
			encrypted, err := crypto.Encrypt([]byte(cfg.Servers[i].Password), newKey)
			if err != nil {
				return err
			}
			cfg.Servers[i].Password = encrypted
		}
	}

	// Update salt and hash
	cfg.Salt = crypto.Base64Encode(newSalt)
	cfg.MasterPasswordHash = crypto.HashPassword([]byte(newPassword), newSalt)

	s.masterKey = newKey
	s.isUnlocked = true

	return s.SaveRaw(cfg)
}

func (s *Store) Load() (*config.StoresConfig, error) {
	cfg, err := s.LoadRaw()
	if err != nil {
		return nil, err
	}

	// If no master password is set, return as-is (no encryption)
	if cfg.Salt == "" || !s.isUnlocked {
		return cfg, nil
	}

	// Decrypt passwords
	for i := range cfg.Servers {
		if cfg.Servers[i].Password != "" {
			decrypted, err := crypto.Decrypt(cfg.Servers[i].Password, s.masterKey)
			if err != nil {
				// If decryption fails, keep the encrypted value
				continue
			}
			cfg.Servers[i].Password = decrypted
		}
	}

	return cfg, nil
}

func (s *Store) Save(cfg *config.StoresConfig) error {
	if !s.isUnlocked || cfg.Salt == "" {
		// No master password, save as-is
		return s.SaveRaw(cfg)
	}

	// Create a copy to encrypt
	cfgCopy := *cfg
	cfgCopy.Servers = make([]config.ServerConfig, len(cfg.Servers))
	for i := range cfg.Servers {
		cfgCopy.Servers[i] = cfg.Servers[i]
		if cfg.Servers[i].Password != "" {
			encrypted, err := crypto.Encrypt([]byte(cfg.Servers[i].Password), s.masterKey)
			if err != nil {
				return err
			}
			cfgCopy.Servers[i].Password = encrypted
		}
	}

	return s.SaveRaw(&cfgCopy)
}

func (s *Store) LoadRaw() (*config.StoresConfig, error) {
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

func (s *Store) SaveRaw(cfg *config.StoresConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}
