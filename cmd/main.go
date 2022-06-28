package main

import (
	"fmt"
	"os"

	"github.com/closmarfer/glacier-backup/pkg/backup/handlers"
)

func main() {
	if len(os.Args) == 1 {
		handler := handlers.NewHandler()
		handler.Run()
		return
	}

	arg := os.Args[1]
	switch arg {
	case "--sizeCount":
		handler := handlers.NewSizeCounter()
		handler.Run()
	default:
		fmt.Printf("Error: invalid argument '%v' provided\n", arg)
		os.Exit(1)
	}
}
