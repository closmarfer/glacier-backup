package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/handlers"
	"github.com/closmarfer/glacier-backup/pkg/backup/serviceprovider"
)

const databaseName = "backup.db"

const appVersion = "3.0"

func main() {
	var selectedRemote string

	if len(os.Args) == 1 {
		printHelp()
	}

	if len(os.Args) != 3 {
		fmt.Println("Error: invalid number of options")
		printHelp()
		os.Exit(1)
	}

	selectedRemote = os.Args[1]
	cfg, err := serviceprovider.ProvideBackupConfiguration(selectedRemote)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	repo, err := serviceprovider.ProvideRemoteFilesRepository(cfg)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	eChecker := backup.NewSQLiteChecker(backup.SqliteConfig{
		Path: cfg.GlacierPath + string(os.PathSeparator) + databaseName,
		Key:  databaseName,
	}, repo)

	arg := os.Args[2]
	switch arg {
	case "--backup":
		back := backup.NewBackuper(repo, eChecker, cfg)
		handler := handlers.NewHandler(eChecker, back)
		handler.Run()
	case "--sizeCount":
		handler := handlers.NewSizeCounter(eChecker)
		handler.Run()
	case "--cleanRemote":
		handler := handlers.NewRemoteCleaner(eChecker, repo)
		handler.Run()
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	help := "glacier-backup [remote] [--sizeCount] [--cleanRemote] [--backup]"
	fmt.Printf("Application version %v %v/%v\n", appVersion, runtime.GOOS, runtime.GOARCH)
	fmt.Println("Usage: " + help)
}
