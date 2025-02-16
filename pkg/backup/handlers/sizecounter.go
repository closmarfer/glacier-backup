package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/closmarfer/glacier-backup/pkg/backup"
)

type SizeCounter struct {
	cfg     backup.Config
	repo    backup.RemoteFilesRepository
	checker backup.ExistentFilesChecker
}

func NewSizeCounter(checker backup.ExistentFilesChecker) SizeCounter {
	return SizeCounter{checker: checker}
}

func (s SizeCounter) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.checker.Open(ctx)
	if err != nil {
		fmt.Printf("Error opening: %v\n", err.Error())
		return
	}

	size := int64(0)
	for path := range s.checker.GetFiles() {
		finfo, err := os.Stat(path)
		if err != nil {
			continue
		}
		size += finfo.Size()
	}

	fmt.Printf("Size of uploaded files: %v GB\n", size/1000_000_000)
}
