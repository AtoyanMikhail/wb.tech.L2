package store

import (
	"os"
	"path/filepath"
)

// Save writes data to file, creating parent dirs
func Save(filename string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	tmp := filename + ".part"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filename)
}
