package backup

import (
	"fmt"
	"os"
	"strings"
)

const glacierBackupFolder = ".glacier-backup"

type Config struct {
	PathsToBackup   []string
	IgnoredPatterns []string
	SelectedRemote  string
	GlacierPath     string
}

const (
	pathsToBackupKey   = "GLACIER_BACKUP_PATHS_TO_BACKUP"
	ignoredPatternsKey = "GLACIER_BACKUP_IGNORED_PATTERNS"
)

var requiredVariables = []string{pathsToBackupKey, ignoredPatternsKey}

type ConfigDecoder struct {
}

func NewConfigDecoder() *ConfigDecoder {
	return &ConfigDecoder{}
}

func (cd ConfigDecoder) LoadConfiguration(selectedRemote string) (Config, error) {
	for _, s := range requiredVariables {
		if os.Getenv(s) == "" {
			return Config{}, fmt.Errorf("environment variable %s not set", s)
		}
	}
	pathsToBackup := strings.Split(os.Getenv(pathsToBackupKey), ";")
	ignoredPatterns := strings.Split(os.Getenv(ignoredPatternsKey), ";")

	userHome, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		PathsToBackup:   pathsToBackup,
		IgnoredPatterns: ignoredPatterns,
		SelectedRemote:  selectedRemote,
		GlacierPath:     userHome,
	}

	return cfg, nil
}

func (conf Config) getApplicationPath(path string) (string, error) {
	userHome, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	return userHome + sep + glacierBackupFolder + sep + path, nil

}

func (conf Config) IsLocal() bool {
	return conf.SelectedRemote == "local"
}

func (conf Config) IsS3() bool {
	return conf.SelectedRemote == "s3"
}
