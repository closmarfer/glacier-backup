package serviceprovider

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	s3aws "github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Config struct {
	Bucket      string
	region      string
	profileName string
}

func newS3Config(bucket string, region string, profileName string) s3Config {
	return s3Config{Bucket: bucket, region: region, profileName: profileName}
}

func provideS3Client(cfg s3Config) (*s3aws.Client, error) {
	s3Configuration, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.region), config.WithSharedConfigProfile(cfg.profileName))

	if err != nil {
		log.Fatal(err)
	}

	client := s3aws.NewFromConfig(s3Configuration)

	return client, nil
}
