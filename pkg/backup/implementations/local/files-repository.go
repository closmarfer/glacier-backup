package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/closmarfer/glacier-backup/pkg/backup"
)

type repository struct {
	localPath string
	timeout   time.Duration
}

func (r repository) Download(_ context.Context, key string, path string) error {
	file, err := os.ReadFile(key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0600)
}

func (r repository) Delete(_ context.Context, remotePath string) error {
	newPath := r.getPath(remotePath)
	time.Sleep(r.timeout)
	return os.Remove(newPath)
}

func NewRepository(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	localPath, ok := cfg.Remotes["local"].CustomConfig["localPath"]
	if !ok {
		return repository{}, errors.New("localPath not defined for local remote in config.yaml")
	}

	timeout, ok := cfg.Remotes["local"].CustomConfig["timeout"]
	if !ok {
		return repository{}, errors.New("timeout not defined for local remote in config.yaml")
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return repository{}, err
	}

	return repository{localPath: localPath, timeout: duration}, nil
}

func (r repository) PutGlacier(ctx context.Context, localPath string) error {
	return r.put(ctx, localPath, localPath)
}

func (r repository) PutEditable(_ context.Context, localPath string, remotePath string) error {
	fileContents, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	return os.WriteFile(remotePath, fileContents, os.ModePerm)
}

func (r repository) put(_ context.Context, localPath string, remotePath string) error {
	sep := string(os.PathSeparator)
	fileContents, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	newPath := r.getPath(remotePath)

	parts := strings.Split(newPath, sep)

	parts2 := parts[0 : len(parts)-1]

	newPath2 := strings.Join(parts2, sep)

	time.Sleep(r.timeout)

	err = os.MkdirAll(newPath2, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(newPath, fileContents, os.ModePerm)
}

func (r repository) Get(_ context.Context, remotePath string) (string, error) {
	file, err := os.ReadFile(r.getPath(remotePath))
	if err != nil {
		return "", backup.NewFileNotFoundError(remotePath)
	}
	return string(file), err
}

func (r repository) getPath(remotePath string) string {
	return fmt.Sprintf("%v/%v", r.localPath, remotePath)
}
