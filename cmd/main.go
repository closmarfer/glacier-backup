package main

import (
	"fmt"
	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/serviceprovider"
	"os"

	"github.com/closmarfer/glacier-backup/pkg/backup/handlers"
)

func main() {
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

	if len(os.Args) == 1 {
		handler := handlers.NewHandler(eChecker, back)
		handler.Run()
		return
	}
	help := "glacier-backup [--sizeCount] [--cleanRemote]"
	if len(os.Args) > 2 {
		fmt.Println("Error: number of max options: 1", help)
		os.Exit(1)
	}

	arg := os.Args[1]
	switch arg {
	case "--sizeCount":
		handler := handlers.NewSizeCounter(eChecker)
		handler.Run()
	case "--cleanRemote":
		handler := handlers.NewRemoteCleaner(eChecker, repo)
		handler.Run()
	default:
		fmt.Printf("Error: invalid argument '%v' provided\n", arg, help)
		os.Exit(1)
	}
}
