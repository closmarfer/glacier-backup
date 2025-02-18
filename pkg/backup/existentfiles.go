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
	uploaded        int
	ignored         int
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
		h.Add(path, t, 0)
	}

	return nil
}

func (h *Checker) Add(path string, uploadedAt time.Time, _ int64) {
	lock.Lock()
	defer lock.Unlock()
	h.files[path] = uploadedAt
	h.uploaded++
}

func (h *Checker) Remove(path string) {
	lock.Lock()
	defer lock.Unlock()
	delete(h.files, path)
}

func (h *Checker) Exists(path string, lastUpdate time.Time) bool {
	lock.Lock()
	defer lock.Unlock()
	uploaded, ok := h.files[path]

	if !ok {
		return false
	}

	isAfter := uploaded.After(lastUpdate)
	if isAfter {
		h.ignored++
	}

	return isAfter
}

func (h *Checker) Close(ctx context.Context) error {
	uPath, _ := h.cfg.getApplicationPath("uploaded_files.csv")
	csvFile, err := os.Create(uPath)
	if err != nil {
		return err
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			fmt.Println("error closing file: " + err.Error())
		}
	}(csvFile)

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
	return os.Remove(uPath)
}

func (h *Checker) Ignored() int {
	lock.Lock()
	defer lock.Unlock()
	return h.ignored
}

func (h *Checker) Uploaded() int {
	lock.Lock()
	defer lock.Unlock()
	return h.uploaded
}

func (h *Checker) GetFiles() map[string]time.Time {
	return h.files
}
