package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/closmarfer/glacier-backup/pkg/backup"
)

type RemoteCleaner struct {
	repo    backup.RemoteFilesRepository
	checker *backup.Checker
}

func NewRemoteCleaner(checker *backup.Checker, repo backup.RemoteFilesRepository) RemoteCleaner {
	return RemoteCleaner{checker: checker, repo: repo}
}

func (c RemoteCleaner) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.checker.Open(ctx)
	if err != nil {
		fmt.Printf("Error opening: %v\n", err.Error())
		return
	}
	defer func() {
		e := c.checker.Close(ctx)
		if e != nil {
			fmt.Println(err.Error())
		}
	}()
	deletedFiles := 0
	for path := range c.checker.GetFiles() {
		_, err := os.Stat(path)
		if err == nil {
			continue
		}
		err = c.repo.Delete(ctx, path)
		if err != nil {
			fmt.Printf("Error deleting file: %v\n", err.Error())
			continue
		}
		c.checker.Remove(path)
		deletedFiles++
	}

	fmt.Printf("Deleted files from remote repository: %d\n", deletedFiles)
}
