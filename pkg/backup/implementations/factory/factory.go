package factory

import (
	"errors"

	"github.com/closmarfer/glacier-backup/pkg/backup"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/local"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/s3"
)

const (
	localRemote = "local"
	s3Remote    = "s3"
)

type FilesRepositoryFactory struct {
	cfg             backup.Config
	localRepository local.Repository
	s3Repository    s3.Repository
}

func NewFilesRepositoryFactory(
	cfg backup.Config,
	localRepository local.Repository,
	s3Repository s3.Repository,
) *FilesRepositoryFactory {
	return &FilesRepositoryFactory{
		cfg:             cfg,
		s3Repository:    s3Repository,
		localRepository: localRepository,
	}
}

func (f FilesRepositoryFactory) Make(remote string) (backup.RemoteFilesRepository, error) {
	if remote == localRemote {
		return f.localRepository, nil
	}

	if remote == s3Remote {

		return f.s3Repository, nil
	}

	return nil, errors.New("Invalid remote")
}
