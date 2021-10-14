package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"glacier-backup/pkg/backup"
	"glacier-backup/pkg/backup/serviceprovider"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Run() {
	cfg, err := serviceprovider.ProvideBackupConfiguration()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	repo, err := serviceprovider.ProvideRemoteFilesRepository(cfg)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	eChecker := backup.NewChecker(repo, cfg)

	back := backup.NewBackuper(repo, eChecker, cfg)

	ctx, cancel := context.WithCancel(context.Background())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	errChan := make(chan error)

	go func() {
		defer cancel()
		errChan <- back.Upload(ctx, errChan)
	}()

	for {
		select {
		case <-shutdown:
			fmt.Println("Stopping program")
			err := eChecker.Close(ctx)
			if err != nil {
				fmt.Println(fmt.Sprintf("error closing existent files checker: %v", err))
			}
			cancel()
		case <-ctx.Done():
			fmt.Println(fmt.Sprintf("Context finished"))
			return
		case err := <-errChan:
			if err != nil {
				fmt.Println(fmt.Sprintf("Error processing files: %v", err))
			}

		}
	}
}
