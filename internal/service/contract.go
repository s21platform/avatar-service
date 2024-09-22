package service

import "context"

type DBRepo interface {
	SetAvatar(userUUID, link string) error
	GetAllAvatars(userUUID string) ([]string, error)
}

type S3Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) (string, error)
}
