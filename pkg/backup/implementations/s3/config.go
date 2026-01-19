package s3

import (
	"fmt"
	"os"
)

type Config struct {
	Bucket      string
	Region      string
	ProfileName string
}

const (
	bucketsKey = "GLACIER_BACKUP_S3_BUCKETS"
	regionKey  = "GLACIER_BACKUP_S3_REGION"
	profileKey = "GLACIER_BACKUP_S3_PROFILE"
)

var requiredVariables = []string{bucketsKey, regionKey, profileKey}

func NewConfig() (Config, error) {
	for _, v := range requiredVariables {
		if os.Getenv(v) == "" {
			return Config{}, fmt.Errorf("environment variable %s must be set", v)
		}
	}
	return Config{
		Bucket:      os.Getenv(bucketsKey),
		Region:      os.Getenv(regionKey),
		ProfileName: os.Getenv(profileKey),
	}, nil
}
