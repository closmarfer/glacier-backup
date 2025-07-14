package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	http2 "net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/closmarfer/glacier-backup/pkg/backup"
)

type Config struct {
	Bucket string
}

type repository struct {
	config Config
	client *s3.Client
}

func (repo repository) Delete(ctx context.Context, remotePath string) error {
	if string(remotePath[0]) == "/" {
		remotePath = remotePath[1:]
	}

	_, err := repo.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(repo.config.Bucket),
		Key:    aws.String(remotePath),
	})

	return err
}

func NewS3Repository(config Config, client *s3.Client) backup.RemoteFilesRepository {
	return repository{config: config, client: client}
}

func (repo repository) PutGlacier(ctx context.Context, localPath string) error {
	return repo.put(ctx, localPath, localPath, types.StorageClassDeepArchive)
}

func (repo repository) PutEditable(ctx context.Context, localPath string, remotePath string) error {
	return repo.put(ctx, localPath, remotePath, types.StorageClassStandard)
}

func (repo repository) put(ctx context.Context, localPath string, remotePath string, s types.StorageClass) error {
	if string(remotePath[0]) == "/" {
		remotePath = remotePath[1:]
	}

	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	_, err = repo.client.PutObject(ctx, &s3.PutObjectInput{
		ACL:          types.ObjectCannedACLPrivate,
		Bucket:       aws.String(repo.config.Bucket),
		Key:          aws.String(remotePath),
		Body:         bytes.NewReader(content),
		StorageClass: s,
	})

	return err
}

func (repo repository) Get(ctx context.Context, remotePath string) (string, error) {
	buff, err := repo.getBuffer(ctx, remotePath)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (repo repository) Download(ctx context.Context, key string, path string) error {
	buff, err := repo.getBuffer(ctx, key)
	if err != nil {
		return err
	}
	return os.WriteFile(path, buff.Bytes(), 0644)
}

func (repo repository) getBuffer(ctx context.Context, remotePath string) (*bytes.Buffer, error) {
	object, err := repo.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(repo.config.Bucket),
		Key:    aws.String(remotePath),
	})
	buff := new(bytes.Buffer)
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http2.StatusNotFound {
			return buff, backup.NewFileNotFoundError(remotePath)
		}
		return buff, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing body: " + err.Error())
		}
	}(object.Body)

	_, err = buff.ReadFrom(object.Body)
	return buff, err
}
