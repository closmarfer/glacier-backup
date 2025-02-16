package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/closmarfer/glacier-backup/pkg/backup"
)

const processingFormat = "Processing. Uploaded: %v, Ignored: %v\n"

type Handler struct {
	checker  backup.ExistentFilesChecker
	backuper backup.Backuper
}

func NewHandler(checker backup.ExistentFilesChecker, backuper backup.Backuper) Handler {
	return Handler{checker: checker, backuper: backuper}
}

func (h Handler) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	errChan := make(chan error)

	go func() {
		defer cancel()
		errChan <- h.backuper.Upload(ctx, errChan)
	}()

	tick := time.NewTicker(60 * time.Second)

	for {
		select {
		case t := <-tick.C:
			fmt.Printf("Processing. Time: %v\n", t.Format("2006-02-01 15:04"))
			fmt.Printf(processingFormat, h.checker.Uploaded(), h.checker.Ignored())
		case <-shutdown:
			fmt.Println("Stopping program")
			err := h.checker.Close(ctx)
			if err != nil {
				fmt.Printf("error closing existent files checker: %v", err)
			}
			fmt.Printf(processingFormat, h.checker.Uploaded(), h.checker.Ignored())
			cancel()
		case <-ctx.Done():
			fmt.Printf(processingFormat, h.checker.Uploaded(), h.checker.Ignored())
			fmt.Println("Process finished")
			return
		case err := <-errChan:
			if err != nil {
				fmt.Printf("Error processing files: %v", err)
			}
		}
	}
}
