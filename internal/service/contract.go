package service

import (
	modelAvatar "avatar_service/internal/model/avatar"
	"context"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type DBRepo interface {
	SetAvatar(userUUID, link string) error
	GetAllAvatars(userUUID string) ([]*avatarproto.Avatar, error)
	GetAvatarData(avatarID int) (*modelAvatar.Info, error)
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
