package main


import (
	"context"
	"fmt"
	"os"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/serviceprovider"
)


func main(){
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eChecker.Open(ctx)

	size := int64(0)
	for path, _ := range eChecker.GetFiles() {
		finfo, err := os.Stat(path)
		if err != nil {
			continue 
		}
		size += finfo.Size()
	}

	fmt.Printf("Size of uploaded files: %v GB\n", (size / 1000_000_000))
}