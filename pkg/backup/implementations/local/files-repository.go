package local

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/closmarfer/glacier-backup/pkg/backup"
)

type repository struct {
	destinationPath string
	separator       string
}

func NewRepository(cfg Config) (backup.RemoteFilesRepository, error) {
	return repository{destinationPath: cfg.DestinationPath, separator: string(os.PathSeparator)}, nil
}

func (r repository) Download(_ context.Context, key string, path string) error {
	file, err := os.ReadFile(r.destinationPath + r.separator + key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0600)
}

func (r repository) Delete(_ context.Context, remotePath string) error {
	newPath := r.getPath(remotePath)
	return os.Remove(newPath)
}

func (r repository) PutGlacier(ctx context.Context, localPath string) error {
	return r.put(ctx, localPath, localPath)
}

func (r repository) PutEditable(_ context.Context, localPath string, remotePath string) error {
	fileContents, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	return os.WriteFile(r.destinationPath+r.separator+remotePath, fileContents, os.ModePerm)
}

func (r repository) put(_ context.Context, localPath string, remotePath string) error {
	fileContents, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}
	newPath := r.getPath(remotePath)

	parts := strings.Split(newPath, r.separator)

	parts2 := parts[0 : len(parts)-1]

	newPath2 := strings.Join(parts2, r.separator)

	err = os.MkdirAll(newPath2, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(newPath, fileContents, os.ModePerm)
}

// cleanPath it prevents error on Windows systems
func (r repository) cleanPath(path string) string {
	return strings.Replace(path, ":", r.separator, 1)
}

func (r repository) Get(_ context.Context, remotePath string) (string, error) {
	file, err := os.ReadFile(r.getPath(remotePath))
	if err != nil {
		return "", backup.NewFileNotFoundError(remotePath)
	}
	return string(file), err
}

func (r repository) getPath(remotePath string) string {
	return fmt.Sprintf("%v%v%v", r.destinationPath, r.separator, r.cleanPath(remotePath))
}
