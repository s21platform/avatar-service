package s3

import (
	"avatar_service/internal/config"
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg" // Регистрация формата jpeg для декодирования изображения в convertToWebP()
	_ "image/png"  // Регистрация формата png для декодирования изображения в convertToWebP()

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	MinioClient *minio.Client
}

func New(cfg *config.Config) (*Client, error) {
	minioClient, err := createMinioClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{MinioClient: minioClient}, nil
}

func createMinioClient(cfg *config.Config) (*minio.Client, error) {
	endpoint := cfg.S3Storage.Endpoint
	accessKeyID := cfg.S3Storage.AccessKeyID
	secretAccessKey := cfg.S3Storage.SecretAccessKey

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error minio.New: %w", err)
	}

	return minioClient, nil
}

func (c *Client) UploadFile(ctx context.Context, bucketName, objectName string, imageData []byte,
	contentType string) (string, error) {
	webpImage, err := convertToWebP(imageData)
	if err != nil {
		return "", err
	}

	err = c.uploadToS3(ctx, bucketName, objectName, webpImage, contentType)
	if err != nil {
		return "", err
	}

	return c.generateLink(bucketName, objectName), nil
}

func convertToWebP(imageData []byte) ([]byte, error) {
	img, err := decodeImage(imageData)
	if err != nil {
		return nil, err
	}

	webpData, err := encodeToWebP(img)
	if err != nil {
		return nil, err
	}

	return webpData, nil
}

func decodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error image.Decode: %w", err)
	}

	return img, nil
}

func encodeToWebP(img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 100)

	if err != nil {
		return nil, fmt.Errorf("error encoder.NewLossyEncoderOptions: %w", err)
	}

	err = webp.Encode(&buf, img, options)
	if err != nil {
		return nil, fmt.Errorf("error webp.Encode: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *Client) uploadToS3(ctx context.Context, bucketName, objectName string, imageData []byte,
	contentType string) error {
	reader := bytes.NewReader(imageData)
	imageSize := int64(len(imageData))

	_, err := c.MinioClient.PutObject(ctx, bucketName, objectName, reader, imageSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("error c.MinioClient.PutObject: %w", err)
	}

	return nil
}

func (c *Client) generateLink(bucketName, objectName string) string {
	return fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", bucketName, objectName)
}
