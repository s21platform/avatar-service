package s3storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) (string, error)
}

type Client struct {
	MinioClient *minio.Client
}

func New(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("error minio.New: %w", err)
	}

	return &Client{MinioClient: minioClient}, nil
}

func (c *Client) UploadFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	imageData []byte,
	contentType string,
) (string, error) {
	_, err := c.MinioClient.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return "", fmt.Errorf("error minioClient.PutObject: %w", err)
	}

	link := fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", bucketName, objectName)

	return link, nil
}
