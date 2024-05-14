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

type Repository struct {
	localPath string
	timeout   time.Duration
}

func (r Repository) Delete(_ context.Context, remotePath string) error {
	newPath := r.getPath(remotePath)
	time.Sleep(r.timeout)
	return os.Remove(newPath)
}

func NewRepository(cfg backup.Config) (Repository, error) {
	localPath, ok := cfg.Remotes["local"].CustomConfig["localPath"]
	if !ok {
		return Repository{}, errors.New("localPath not defined for local remote in config.yaml")
	}

	timeout, ok := cfg.Remotes["local"].CustomConfig["timeout"]
	if !ok {
		return Repository{}, errors.New("timeout not defined for local remote in config.yaml")
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return Repository{}, err
	}

	return Repository{localPath: localPath, timeout: duration}, nil
}

func (r Repository) PutGlacier(ctx context.Context, localPath string) error {
	return r.put(ctx, localPath, localPath)
}

func (r Repository) PutEditable(ctx context.Context, localPath string, remotePath string) error {
	return r.put(ctx, localPath, remotePath)
}

func (r Repository) put(_ context.Context, localPath string, remotePath string) error {
	fileContents, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	newPath := r.getPath(remotePath)

	parts := strings.Split(newPath, "/")

	parts2 := parts[0 : len(parts)-1]

	newPath2 := strings.Join(parts2, "/")

	time.Sleep(r.timeout)

	err = os.MkdirAll(newPath2, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(newPath, fileContents, os.ModePerm)
}

func (r Repository) Get(_ context.Context, remotePath string) (string, error) {
	file, err := os.ReadFile(r.getPath(remotePath))
	if err != nil {
		return "", backup.NewFileNotFoundError(remotePath)
	}
	return string(file), err
}

func (r Repository) getPath(remotePath string) string {
	return fmt.Sprintf("%v/%v", r.localPath, remotePath)
}
