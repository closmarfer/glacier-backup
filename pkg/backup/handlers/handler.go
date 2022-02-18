package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/serviceprovider"
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

	tick := time.NewTimer(60 * time.Second)

	for {
		select {
		case t := <-tick.C:
			fmt.Printf("Processing. Time: %v\n", t.Format("2006-02-01 15:04"))
			fmt.Printf("Processing. Uploaded: %v, Ignored: %v\n", eChecker.Uploaded(), eChecker.Ignored())
		case <-shutdown:
			fmt.Println("Stopping program")
			err := eChecker.Close(ctx)
			if err != nil {
				fmt.Printf("error closing existent files checker: %v", err)
			}
			fmt.Printf("Processing. Uploaded: %v, Ignored: %v\n", eChecker.Uploaded(), eChecker.Ignored())
			cancel()
		case <-ctx.Done():
			fmt.Printf("Processing. Uploaded: %v, Ignored: %v\n", eChecker.Uploaded(), eChecker.Ignored())
			fmt.Println("Process finished")
			return
		case err := <-errChan:
			if err != nil {
				fmt.Printf("Error processing files: %v", err)
			}
		}
	}
}
