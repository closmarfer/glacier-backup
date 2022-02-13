package backup

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var lock = sync.RWMutex{}

const timeLayout = "2006-01-02 15:04:05"
const uploadedFilesPath = "uploaded_files.csv"

type Checker struct {
	filesRepository RemoteFilesRepository
	files           map[string]time.Time
	cfg             Config
}

func NewChecker(filesRepository RemoteFilesRepository, cfg Config) *Checker {
	return &Checker{filesRepository: filesRepository, files: map[string]time.Time{}, cfg: cfg}
}

func (h Checker) Open(ctx context.Context) error {
	contents, err := h.filesRepository.Get(ctx, uploadedFilesPath)
	if err != nil {
		if _, ok := err.(FileNotFoundError); !ok {
			return err
		}
	}

	if contents == "" {
		return nil
	}

	r := csv.NewReader(strings.NewReader(contents))

	for {
		cols, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		path := cols[0]
		uploadedAt := cols[1]

		t, err := time.Parse(timeLayout, uploadedAt)
		if err != nil {
			return fmt.Errorf("error parsing dateTime: %w", err)
		}
		h.Add(path, t)
	}

	return nil
}

func (h Checker) Add(path string, uploadedAt time.Time) {
	lock.Lock()
	defer lock.Unlock()
	h.files[path] = uploadedAt
}

func (h Checker) Exists(path string, lastUpdate time.Time) bool {
	lock.Lock()
	defer lock.Unlock()
	uploaded, ok := h.files[path]

	if !ok {
		return false
	}

	return uploaded.After(lastUpdate)
}

func (h *Checker) Close(ctx context.Context) error {
	uPath, _ := h.cfg.getApplicationPath("uploaded_files.csv")
	csvFile, err := os.Open(uPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)

	for path, t := range h.files {
		_ = w.Write([]string{path, t.Format(timeLayout)})
	}
	w.Flush()

	err = h.filesRepository.PutEditable(ctx, uPath, "/uploaded_files.csv")
	if err != nil {
		return fmt.Errorf("error putting uploaded_files.csv file: %w", err)
	}
	h.files = map[string]time.Time{}

	return nil
}
