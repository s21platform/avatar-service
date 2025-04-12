package s3

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg" // Регистрация формата jpeg для декодирования изображения в convertToWebP()
	_ "image/png"  // Регистрация формата png для декодирования изображения в convertToWebP()
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
)

type Client struct {
	minio      *minio.Client
	bucketName string
}

func New(cfg *config.Config) *Client {
	client, err := minio.New(cfg.S3Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3Storage.AccessKeyID, cfg.S3Storage.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal("failed to create S3 client: ", err)
	}

	if cfg.S3Storage.BucketName == "" {
		log.Fatal("bucket name is required")
	}

	return &Client{
		minio:      client,
		bucketName: cfg.S3Storage.BucketName,
	}
}

func (c *Client) PutObject(ctx context.Context, avatar *model.AvatarContent) (string, error) {
	webpImage, err := convertToWebP(avatar.ImageData)
	if err != nil {
		return "", fmt.Errorf("failed to convert to webp: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	fileNameNoExt := strings.TrimSuffix(avatar.Filename, filepath.Ext(avatar.Filename))
	objectName := fmt.Sprintf("%s/%s/%s_%s.webp", avatar.AvatarType, avatar.UUID, timestamp, fileNameNoExt)

	reader := bytes.NewReader(webpImage)
	_, err = c.minio.PutObject(
		ctx,
		c.bucketName,
		objectName,
		reader,
		int64(len(webpImage)),
		minio.PutObjectOptions{ContentType: "image/webp"},
	)
	if err != nil {
		return "", fmt.Errorf("failed to put object: %w", err)
	}

	return fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", c.bucketName, objectName), nil
}

func (c *Client) RemoveObject(ctx context.Context, link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) < 2 {
		return fmt.Errorf("failed to link missing bucket or object name")
	}
	bucketName, objectName := parts[0], parts[1]

	_, err = c.minio.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to check object existence: %w", err)
	}

	err = c.minio.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove object: %w", err)
	}

	return nil
}

func convertToWebP(imageData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder options: %w", err)
	}

	if err = webp.Encode(&buf, img, options); err != nil {
		return nil, fmt.Errorf("failed to encode to WebP: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *Client) BucketName() string {
	return c.bucketName
}
