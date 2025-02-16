package backup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type pathInfo struct {
	path       string
	lastUpdate time.Time
}

type FileNotFoundError struct {
	filePath string
}

func NewFileNotFoundError(filePath string) FileNotFoundError {
	return FileNotFoundError{filePath: filePath}
}

func (f FileNotFoundError) Error() string {
	return "file not found " + f.filePath
}

type RemoteFilesRepository interface {
	PutGlacier(ctx context.Context, localPath string) error
	PutEditable(ctx context.Context, localPath string, remotePath string) error
	Delete(ctx context.Context, remotePath string) error
	Get(ctx context.Context, remotePath string) (string, error)
	Download(ctx context.Context, key string, path string) error
}

type ExistentFilesChecker interface {
	Open(ctx context.Context) error
	Add(path string, uploadedAt time.Time)
	Remove(path string)
	Exists(path string, lastUpdated time.Time) bool
	Close(ctx context.Context) error
	Ignored() int
	Uploaded() int
	GetFiles() map[string]time.Time
}

type Backuper struct {
	filesRepository RemoteFilesRepository
	eChecker        ExistentFilesChecker
	config          Config
}

func NewBackuper(
	filesRepository RemoteFilesRepository,
	eChecker ExistentFilesChecker,
	config Config,
) Backuper {
	return Backuper{
		filesRepository: filesRepository,
		eChecker:        eChecker,
		config:          config,
	}
}

func (h Backuper) Upload(ctx context.Context, errChan chan error) error {
	defer func(eChecker ExistentFilesChecker, ctx context.Context) {
		err := eChecker.Close(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("error closing existentFilesChecker: %v", err))
		}
	}(h.eChecker, ctx)

	var wg sync.WaitGroup
	err := h.eChecker.Open(ctx)
	if err != nil {
		return err
	}

	paths := make(chan pathInfo)

	for i := 0; i < 5; i++ {
		w := newWorker(&wg, h.filesRepository, h.eChecker)

		w.run(ctx, paths, errChan)
	}

	go func() {
		defer close(paths)
		for _, p := range h.config.PathsToBackup {
			err := h.iterate(ctx, p, paths)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}
	}()

	wg.Wait()

	return nil
}

func (h Backuper) iterate(ctx context.Context, path string, paths chan<- pathInfo) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return io.EOF
		default:
			if err != nil {
				return fmt.Errorf("error walking file %s: %w", path, err)
			}

			if info.IsDir() || h.shouldBeIgnored(path) {
				return nil
			}

			paths <- pathInfo{
				path:       path,
				lastUpdate: info.ModTime().UTC(),
			}

			return nil
		}

	})

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return fmt.Errorf("error iterating folder: %w", err)
	}
	return nil
}

func (h Backuper) shouldBeIgnored(path string) bool {
	for _, ignoredPattern := range h.config.IgnoredPatterns {
		r := regexp.MustCompile(ignoredPattern)
		if r.MatchString(path) {
			return true
		}
	}

	return false
}

type worker struct {
	wg              *sync.WaitGroup
	filesRepository RemoteFilesRepository
	existent        ExistentFilesChecker
}

func newWorker(wg *sync.WaitGroup, filesRepository RemoteFilesRepository, existent ExistentFilesChecker) *worker {
	return &worker{wg: wg, filesRepository: filesRepository, existent: existent}
}

func (w worker) run(ctx context.Context, paths <-chan pathInfo, errChan chan error) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case path, ok := <-paths:
				if !ok {
					return
				}
				if w.existent.Exists(path.path, path.lastUpdate) {
					break
				}
				err := w.filesRepository.PutGlacier(ctx, path.path)
				if err != nil {
					errChan <- fmt.Errorf("error putting file: %w", err)
				}
				w.existent.Add(path.path, time.Now().UTC())
			case <-ctx.Done():
				return
			}
		}
	}()
}
