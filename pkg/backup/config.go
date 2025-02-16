package backup

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const glacierBackupFolder = ".glacier-backup"

type Remote struct {
	CustomConfig map[string]string `yaml:"customConfig"`
}

type Config struct {
	PathsToBackup   []string          `yaml:"pathsToBackup"`
	IgnoredPatterns []string          `yaml:"ignoredPatterns"`
	SelectedRemote  string            `yaml:"selectedRemote"`
	DatabaseKey     string            `yaml:"databaseKey" default:"backup.db"`
	Remotes         map[string]Remote `yaml:"remotes"`
	GlacierPath     string
}

type ConfigDecoder struct {
}

func NewConfigDecoder() *ConfigDecoder {
	return &ConfigDecoder{}
}

func (cd ConfigDecoder) LoadConfiguration() (Config, error) {
	var cfg Config
	userHome, err := os.UserHomeDir()

	if err != nil {
		return cfg, err
	}

	ypath, err := cfg.getApplicationPath("config.yaml")

	if err != nil {
		return cfg, err
	}

	yfile, err := ioutil.ReadFile(ypath)

	if err != nil {
		return cfg, fmt.Errorf("error decoding yaml: %w", err)
	}

	err = yaml.Unmarshal(yfile, &cfg)

	if err != nil {
		return cfg, fmt.Errorf("error unmarshal yaml: %w", err)
	}

	cfg.GlacierPath = userHome

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
