package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
	"github.com/s21platform/avatar-service/pkg/avatar"
	"github.com/s21platform/avatar-service/pkg/new_avatar_register"
)

type Service struct {
	avatar.UnimplementedAvatarServiceServer
	s3Client             S3Storage
	repository           DBRepo
	userKafkaProducer    KafkaProducer
	societyKafkaProducer KafkaProducer
}

func New(s3Client S3Storage, repo DBRepo, userKafkaProducer KafkaProducer, societyKafkaProducer KafkaProducer) *Service {
	return &Service{
		s3Client:             s3Client,
		repository:           repo,
		userKafkaProducer:    userKafkaProducer,
		societyKafkaProducer: societyKafkaProducer,
	}
}

func (s *Service) SetUserAvatar(stream avatar.AvatarService_SetUserAvatarServer) error {
	ctx := stream.Context()
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetUserAvatar")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("uuid is required")
		return status.Error(codes.InvalidArgument, "uuid is required")
	}

	avatarContent := &model.AvatarContent{
		AvatarType: model.UserAvatarType,
		UUID:       uuid,
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logger.Error(fmt.Sprintf("failed to receive data from stream: %v", err))
			return status.Errorf(codes.Internal, "failed to receive data from stream: %v", err)
		}

		if avatarContent.Filename == "" {
			avatarContent.Filename = in.Filename
		}

		avatarContent.ImageData = append(avatarContent.ImageData, in.Batch...)
	}

	link, err := s.s3Client.PutObject(ctx, avatarContent)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to upload file to S3: %v", err))
		return status.Errorf(codes.Internal, "failed to upload file to S3: %v", err)
	}

	if err = s.repository.SetUserAvatar(ctx, uuid, link); err != nil {
		logger.Error(fmt.Sprintf("failed to save avatar to database: %v", err))
		return status.Errorf(codes.Internal, "failed to save avatar to database: %v", err)
	}

	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: uuid,
		Link: link,
	}
	if err = s.userKafkaProducer.ProduceMessage(ctx, msg, uuid); err != nil {
		logger.Error(fmt.Sprintf("failed to produce message to user service: %v", err))
		return status.Errorf(codes.Internal, "failed to produce message to user service: %v", err)
	}

	return stream.SendAndClose(&avatar.SetUserAvatarOut{
		Link: link,
	})
}

func (s *Service) GetAllUserAvatars(ctx context.Context, _ *emptypb.Empty) (*avatar.GetAllUserAvatarsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAllUserAvatars")

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("uuid is required")
		return nil, status.Error(codes.InvalidArgument, "uuid is required")
	}

	avatars, err := s.repository.GetAllUserAvatars(ctx, uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all user avatars: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get all user avatars: %v", err)
	}

	return &avatar.GetAllUserAvatarsOut{
		AvatarList: avatars.FromDTO(),
	}, nil
}

func (s *Service) DeleteUserAvatar(ctx context.Context, in *avatar.DeleteUserAvatarIn) (*avatar.Avatar, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("DeleteUserAvatar")

	avatarInfo, err := s.repository.GetUserAvatarData(ctx, int(in.AvatarId))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user avatar: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get avatar data: %v", err)
	}

	err = s.s3Client.RemoveObject(ctx, avatarInfo.Link)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete user avatar: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete avatar in s3: %v", err)
	}

	err = s.repository.DeleteUserAvatar(ctx, avatarInfo.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete user avatar: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete avatar in postgres: %v", err)
	}

	latestAvatar := s.repository.GetLatestUserAvatar(ctx, avatarInfo.UUID)

	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: avatarInfo.UUID,
		Link: latestAvatar,
	}
	err = s.userKafkaProducer.ProduceMessage(ctx, msg, avatarInfo.UUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to produce avatar to user service: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to produce avatar to user service: %v", err)
	}

	return &avatar.Avatar{
		Id:   int32(avatarInfo.ID),
		Link: avatarInfo.Link,
	}, nil
}

func (s *Service) SetSocietyAvatar(stream avatar.AvatarService_SetSocietyAvatarServer) error {
	ctx := stream.Context()
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetSocietyAvatar")

	avatarContent := &model.AvatarContent{AvatarType: model.SocietyAvatarType}

	for {
		in, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logger.Error(fmt.Sprintf("failed to receive data from stream: %v", err))
			return status.Errorf(codes.Internal, "failed to receive data from stream: %v", err)
		}

		if avatarContent.UUID == "" || avatarContent.Filename == "" {
			avatarContent.UUID = in.Uuid
			avatarContent.Filename = in.Filename
		}

		avatarContent.ImageData = append(avatarContent.ImageData, in.Batch...)
	}

	if avatarContent.UUID == "" {
		logger.Error("society uuid is required")
		return status.Error(codes.InvalidArgument, "society uuid is required")
	}

	link, err := s.s3Client.PutObject(ctx, avatarContent)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to upload file to S3: %v", err))
		return status.Errorf(codes.Internal, "failed to upload file to S3: %v", err)
	}

	if err = s.repository.SetSocietyAvatar(ctx, avatarContent.UUID, link); err != nil {
		logger.Error(fmt.Sprintf("failed to save avatar to database: %v", err))
		return status.Errorf(codes.Internal, "failed to save avatar to database: %v", err)
	}

	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: avatarContent.UUID,
		Link: link,
	}
	err = s.societyKafkaProducer.ProduceMessage(ctx, msg, avatarContent.UUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to produce message to society service: %v", err))
		return status.Errorf(codes.Internal, "failed to produce message to society service: %v", err)
	}

	return stream.SendAndClose(&avatar.SetSocietyAvatarOut{
		Link: link,
	})
}

func (s *Service) GetAllSocietyAvatars(ctx context.Context, in *avatar.GetAllSocietyAvatarsIn) (*avatar.GetAllSocietyAvatarsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAllSocietyAvatars")

	if in.Uuid == "" {
		logger.Error("society uuid is required")
		return nil, status.Error(codes.InvalidArgument, "society uuid is required")
	}

	avatars, err := s.repository.GetAllSocietyAvatars(ctx, in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all society avatars: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get all society avatars: %v", err)
	}

	return &avatar.GetAllSocietyAvatarsOut{
		AvatarList: avatars.FromDTO(),
	}, nil
}

func (s *Service) DeleteSocietyAvatar(ctx context.Context, in *avatar.DeleteSocietyAvatarIn) (*avatar.Avatar, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("DeleteSocietyAvatar")

	avatarInfo, err := s.repository.GetSocietyAvatarData(ctx, int(in.AvatarId))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get avatar data: %v", err))
		return nil, status.Errorf(codes.NotFound, "failed to get avatar data: %v", err)
	}

	err = s.s3Client.RemoveObject(ctx, avatarInfo.Link)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete avatar in s3: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete avatar in s3: %v", err)
	}

	err = s.repository.DeleteSocietyAvatar(ctx, avatarInfo.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete avatar in postgres: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete avatar in postgres: %v", err)
	}

	latestAvatar := s.repository.GetLatestSocietyAvatar(ctx, avatarInfo.UUID)

	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: avatarInfo.UUID,
		Link: latestAvatar,
	}
	err = s.societyKafkaProducer.ProduceMessage(ctx, msg, avatarInfo.UUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to produce avatar: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to produce avatar: %v", err)
	}

	return &avatar.Avatar{
		Id:   int32(avatarInfo.ID),
		Link: avatarInfo.Link,
	}, nil
}
