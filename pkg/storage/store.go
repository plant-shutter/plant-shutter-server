package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/plant-shutter/plant-shutter-server/pkg/utils/config"
)

type Storage struct {
	path string
}

func New(cfg config.ImageStorage) (*Storage, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("path can not be empty")
	}
	return &Storage{path: cfg.Path}, nil
}

func (s *Storage) Save(deviceName, project, fileName string, src io.Reader) error {
	dst := s.GetPath(deviceName, project, fileName)
	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *Storage) GetPath(deviceName, project, fileName string) string {
	return path.Join(s.path, deviceName, project, fileName)
}
