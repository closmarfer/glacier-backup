package main

import (
	"glacier-backup/pkg/backup/handlers"
)

func main() {
	handler := handlers.NewHandler()
	handler.Run()
}
