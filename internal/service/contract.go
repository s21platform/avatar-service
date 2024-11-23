package service

import (
	"avatar_service/internal/model"
	"context"
)

type DBRepo interface {
	SetAvatar(userUUID, link string) error
	GetAllAvatars(userUUID string) (*model.AvatarInfoList, error)
	GetAvatarData(avatarID int) (*model.AvatarInfo, error)
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
