package serviceprovider

import (
	"context"
	"github.com/closmarfer/glacier-backup/pkg/backup/implementations/s3"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	s3aws "github.com/aws/aws-sdk-go-v2/service/s3"
)

func provideS3Client(cfg s3.Config) (*s3aws.Client, error) {
	s3Configuration, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region), config.WithSharedConfigProfile(cfg.ProfileName))

	if err != nil {
		log.Fatal(err)
	}

	client := s3aws.NewFromConfig(s3Configuration)

	return client, nil
}
