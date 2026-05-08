package serviceprovider

import (
	"fmt"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/local"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/s3"
)

func ProvideBackupConfiguration(selectedRemote string) (backup.Config, error) {
	cfgDecoder := backup.NewConfigDecoder()
	return cfgDecoder.LoadConfiguration(selectedRemote)
}

func ProvideRemoteFilesRepository(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	return makeRepositoryOrFail(cfg)
}

func makeRepositoryOrFail(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	fmt.Printf("Selected remote: '%v'\n", cfg.SelectedRemote)
	if cfg.IsLocal() {
		c, err := local.NewConfig()
		if err != nil {
			return nil, err
		}

		localRepository, err := local.NewRepository(c)
		if err != nil {
			return nil, err
		}
		return localRepository, nil
	}

	if cfg.IsS3() {
		s3cfg, err := s3.NewConfig()
		if err != nil {
			return nil, err
		}

		client, err := provideS3Client(s3cfg)

		if err != nil {
			return nil, fmt.Errorf("s3 client could not be created: %w", err)
		}

		s3Repository := s3.NewS3Repository(
			s3cfg,
			client,
		)

		return s3Repository, nil
	}

	return nil, fmt.Errorf("not supported remote '%v'", cfg.SelectedRemote)
}
