//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
//go:generate mockgen -destination=mock_stream_test.go -package=${GOPACKAGE} github.com/s21platform/avatar-service/pkg/avatar AvatarService_SetUserAvatarServer,AvatarService_SetSocietyAvatarServer
package service

import (
	"context"

	"github.com/s21platform/avatar-service/internal/model"
)

type DBRepo interface {
	SetUserAvatar(ctx context.Context, UUID, link string) error
	GetAllUserAvatars(ctx context.Context, UUID string) (*model.AvatarMetadataList, error)
	GetUserAvatarData(ctx context.Context, avatarID int) (*model.AvatarMetadata, error)
	DeleteUserAvatar(ctx context.Context, avatarID int) error
	GetLatestUserAvatar(ctx context.Context, UUID string) string

	SetSocietyAvatar(ctx context.Context, UUID, link string) error
	GetAllSocietyAvatars(ctx context.Context, UUID string) (*model.AvatarMetadataList, error)
	GetSocietyAvatarData(ctx context.Context, avatarID int) (*model.AvatarMetadata, error)
	DeleteSocietyAvatar(ctx context.Context, avatarID int) error
	GetLatestSocietyAvatar(ctx context.Context, UUID string) string
}

type S3Storage interface {
	PutObject(ctx context.Context, avatar *model.AvatarContent) (string, error)
	RemoveObject(ctx context.Context, link string) error
	BucketName() string
}

type KafkaProducer interface {
	ProduceMessage(ctx context.Context, message interface{}, key interface{}) error
}
