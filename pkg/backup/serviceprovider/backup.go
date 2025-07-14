package serviceprovider

import (
	"errors"
	"fmt"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/local"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/s3"
)

const s3Key = "s3"

func ProvideBackupConfiguration() (backup.Config, error) {
	cfgDecoder := backup.NewConfigDecoder()
	return cfgDecoder.LoadConfiguration()
}

func ProvideRemoteFilesRepository(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	return makeRepositoryOrFail(cfg)
}

func makeRepositoryOrFail(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	fmt.Printf("Selected remote: '%v'\n", cfg.SelectedRemote)
	if cfg.IsLocal() {
		localRepository, err := local.NewRepository(cfg)
		if err != nil {
			return nil, err
		}
		return localRepository, nil
	}

	if cfg.IsS3() {
		s3Config, ok := cfg.Remotes[s3Key]
		if !ok {
			return nil, errors.New("configuration not defined for remote s3 in config.yaml")
		}

		configUser := getS3Config(s3Config)

		client, err := provideS3Client(configUser)

		if err != nil {
			return nil, fmt.Errorf("s3 client could not be created: %w", err)
		}

		s3Repository := s3.NewS3Repository(
			s3.Config{Bucket: configUser.Bucket},
			client,
		)

		return s3Repository, nil
	}

	return nil, fmt.Errorf("not supported remote %v", cfg.SelectedRemote)
}

func getS3Config(remote backup.Remote) s3Config {
	custom := remote.CustomConfig
	return newS3Config(custom["bucket"], custom["region"], custom["profileName"])
}
