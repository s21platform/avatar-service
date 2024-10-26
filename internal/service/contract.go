package service

import (
	"context"
	"time"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type DBRepo interface {
	SetAvatar(userUUID, link string) error
	GetAllAvatars(userUUID string) ([]*avatarproto.Avatar, error)
	GetAvatarData(avatarID int) (int, string, string, time.Time, error)
	DeleteAvatar(avatarID int) error
	GetLatestAvatar(userUUID string) string
}

type S3Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, link string) error
}

type NewAvatarRegisterSrv interface {
	ProduceMessage(message interface{}) error
}
