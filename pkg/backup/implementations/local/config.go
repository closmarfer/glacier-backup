package local

import (
	"fmt"
	"os"
)

type Config struct {
	DestinationPath string
}

const (
	destinationPathKey = "GLACIER_BACKUP_LOCAL_DESTINATION_PATH"
)

var requiredVariables = []string{destinationPathKey}

func NewConfig() (Config, error) {
	for _, v := range requiredVariables {
		if os.Getenv(v) == "" {
			return Config{}, fmt.Errorf("environment variable %s must be set", v)
		}
	}
	return Config{
		DestinationPath: os.Getenv(destinationPathKey),
	}, nil
}
