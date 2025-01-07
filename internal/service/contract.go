package service

import (
	"avatar_service/internal/model"
	"context"
)

type DBRepo interface {
	SetUserAvatar(UUID, link string) error
	GetAllUserAvatars(UUID string) (*model.AvatarInfoList, error)
	GetUserAvatarData(avatarID int) (*model.AvatarInfo, error)
	DeleteUserAvatar(avatarID int) error
	GetLatestUserAvatar(UUID string) string

	SetSocietyAvatar(UUID, link string) error
	GetAllSocietyAvatars(UUID string) (*model.AvatarInfoList, error)
	GetSocietyAvatarData(avatarID int) (*model.AvatarInfo, error)
	DeleteSocietyAvatar(avatarID int) error
	GetLatestSocietyAvatar(UUID string) string
}

type S3Storage interface {
	UploadFile(ctx context.Context, bucketName, objectName string, fileData []byte, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, link string) error
}

type NewAvatarRegisterSrv interface {
	ProduceMessage(message interface{}) error
}
