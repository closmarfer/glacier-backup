package main

import (
	"github.com/closmarfer/glacier-backup/pkg/backup/handlers"
)

func main() {
	handler := handlers.NewHandler()
	handler.Run()
}
