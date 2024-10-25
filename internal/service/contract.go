package service

import (
	"context"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type DBRepo interface {
	SetAvatar(userUUID, link string) error
	GetAllAvatars(userUUID string) ([]*avatarproto.Avatar, error)
}

type S3Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) (string, error)
}

type NewAvatarRegisterSrv interface {
	ProduceMessage(message interface{}) error
}
