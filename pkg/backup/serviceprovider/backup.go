package serviceprovider

import (
	"errors"
	"fmt"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/factory"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/local"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/s3"
)

const s3Key = "s3"

func ProvideBackupConfiguration() (backup.Config, error) {
	cfgDecoder := backup.NewConfigDecoder()
	return cfgDecoder.LoadConfiguration()
}

func ProvideRemoteFilesRepository(cfg backup.Config) (backup.RemoteFilesRepository, error) {
	fct, err := provideFilesRepositoryFactory(cfg)
	if err != nil {
		return nil, err
	}
	return fct.Make(cfg.SelectedRemote)
}

func provideFilesRepositoryFactory(cfg backup.Config) (*factory.FilesRepositoryFactory, error) {
	localRepository, err := local.NewRepository(cfg)

	if err != nil {
		return nil, err
	}

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

	return factory.NewFilesRepositoryFactory(cfg, localRepository, s3Repository), nil
}

func getS3Config(remote backup.Remote) s3Config {
	custom := remote.CustomConfig
	return newS3Config(custom["bucket"], custom["region"], custom["profileName"])
}
