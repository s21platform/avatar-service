package service

import (
	"avatar_service/internal/config"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/user-proto/user-proto/new_avatar_register"
)

type AvatarType string

const (
	TypeUser    AvatarType = "user"
	TypeSociety AvatarType = "society"
)

type Service struct {
	avatarproto.UnimplementedAvatarServiceServer
	s3Client             S3Storage
	repository           DBRepo
	userKafkaProducer    NewAvatarRegisterSrv
	societyKafkaProducer NewAvatarRegisterSrv
	bucketName           string
}

func New(s3Client S3Storage, repo DBRepo, userKafkaProducer NewAvatarRegisterSrv, societyKafkaProducer NewAvatarRegisterSrv, bucketName string) *Service {
	return &Service{
		s3Client:             s3Client,
		repository:           repo,
		userKafkaProducer:    userKafkaProducer,
		societyKafkaProducer: societyKafkaProducer,
		bucketName:           bucketName,
	}
}

func (s *Service) SetUserAvatar(stream avatarproto.AvatarService_SetUserAvatarServer) error {
	//logger := logger_lib.FromContext(stream.Context(), config.KeyLogger)
	//logger.AddFuncName("SetUserAvatar")

	UUID, filename, imageData, err := s.receiveUserData(stream)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	link, err := s.uploadToS3(UUID, filename, imageData, TypeUser)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	if err = s.repository.SetUserAvatar(UUID, link); err != nil {
		//logger.Error(fmt.Sprintf("failed to save avatar to database: %v", err))
		return fmt.Errorf("failed to save avatar to database: %w", err)
	}

	err = s.produceNewUserAvatar(UUID, link)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	return stream.SendAndClose(&avatarproto.SetUserAvatarOut{
		Link: link,
	})
}

func (s *Service) receiveUserData(stream avatarproto.AvatarService_SetUserAvatarServer) (string, string, []byte, error) {
	var (
		UUID      string
		filename  string
		imageData []byte
	)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", "", nil, fmt.Errorf("failed to receive data from stream: %w", err)
		}

		if UUID == "" && filename == "" {
			UUID = in.Uuid
			filename = in.Filename
		}

		imageData = append(imageData, in.Batch...)
	}

	return UUID, filename, imageData, nil
}

func (s *Service) produceNewUserAvatar(UUID, link string) error {
	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: UUID,
		Link: link,
	}

	err := s.userKafkaProducer.ProduceMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAllUserAvatars(ctx context.Context, in *avatarproto.GetAllUserAvatarsIn) (*avatarproto.GetAllUserAvatarsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAllUserAvatars")

	avatars, err := s.repository.GetAllUserAvatars(in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all user avatars: %v", err))
		return nil, fmt.Errorf("failed to get all user avatars: %w", err)
	}

	return &avatarproto.GetAllUserAvatarsOut{
		AvatarList: avatars.FromDTO(),
	}, nil
}

func (s *Service) DeleteUserAvatar(ctx context.Context, in *avatarproto.DeleteUserAvatarIn) (*avatarproto.Avatar, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("DeleteUserAvatar")

	avatarInfo, err := s.repository.GetUserAvatarData(int(in.AvatarId))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user avatar: %v", err))
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	err = s.s3Client.DeleteAvatar(ctx, avatarInfo.Link)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete user avatar: %v", err))
		return nil, fmt.Errorf("failed to delete avatar in s3: %w", err)
	}

	err = s.repository.DeleteUserAvatar(avatarInfo.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete user avatar: %v", err))
		return nil, fmt.Errorf("failed to delete avatar in postgres: %w", err)
	}

	latestAvatar := s.repository.GetLatestUserAvatar(avatarInfo.UUID)

	err = s.produceNewUserAvatar(avatarInfo.UUID, latestAvatar)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete user avatar: %v", err))
		return nil, fmt.Errorf("failed to produce avatar: %w", err)
	}

	return &avatarproto.Avatar{
		Id:   int32(avatarInfo.ID),
		Link: avatarInfo.Link,
	}, err
}

func (s *Service) SetSocietyAvatar(stream avatarproto.AvatarService_SetSocietyAvatarServer) error {
	//logger := logger_lib.FromContext(stream.Context(), config.KeyLogger)
	//logger.AddFuncName("SetSocietyAvatar")

	UUID, filename, imageData, err := s.receiveSocietyData(stream)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	link, err := s.uploadToS3(UUID, filename, imageData, TypeSociety)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	if err = s.repository.SetSocietyAvatar(UUID, link); err != nil {
		//logger.Error(fmt.Sprintf("failed to save avatar to database: %v", err))
		return fmt.Errorf("failed to save avatar to database: %w", err)
	}

	err = s.produceNewSocietyAvatar(UUID, link)
	if err != nil {
		//logger.Error(fmt.Sprintf("%v", err))
		return err
	}

	return stream.SendAndClose(&avatarproto.SetSocietyAvatarOut{
		Link: link,
	})
}

func (s *Service) receiveSocietyData(stream avatarproto.AvatarService_SetSocietyAvatarServer) (string, string, []byte, error) {
	var (
		UUID      string
		filename  string
		imageData []byte
	)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", "", nil, fmt.Errorf("failed to receive data from stream: %w", err)
		}

		if UUID == "" && filename == "" {
			UUID = in.Uuid
			filename = in.Filename
		}

		imageData = append(imageData, in.Batch...)
	}

	return UUID, filename, imageData, nil
}

func (s *Service) produceNewSocietyAvatar(UUID, link string) error {
	// todo подтянуть контракт из society-proto
	//msg := &new_avatar_register.NewAvatarRegister{
	//	Uuid: UUID,
	//	Link: link,
	//}

	//err := s.societyKafkaProducer.ProduceMessage(msg)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *Service) GetAllSocietyAvatars(ctx context.Context, in *avatarproto.GetAllSocietyAvatarsIn) (*avatarproto.GetAllSocietyAvatarsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAllSocietyAvatars")

	avatars, err := s.repository.GetAllSocietyAvatars(in.Uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all society avatars: %v", err))
		return nil, fmt.Errorf("failed to get all society avatars: %w", err)
	}

	return &avatarproto.GetAllSocietyAvatarsOut{
		AvatarList: avatars.FromDTO(),
	}, nil
}

func (s *Service) DeleteSocietyAvatar(ctx context.Context, in *avatarproto.DeleteSocietyAvatarIn) (*avatarproto.Avatar, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("DeleteSocietyAvatar")

	avatarInfo, err := s.repository.GetSocietyAvatarData(int(in.AvatarId))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get avatar data: %v", err))
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	err = s.s3Client.DeleteAvatar(ctx, avatarInfo.Link)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete avatar in s3: %v", err))
		return nil, fmt.Errorf("failed to delete avatar in s3: %w", err)
	}

	err = s.repository.DeleteSocietyAvatar(avatarInfo.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete avatar in postgres: %v", err))
		return nil, fmt.Errorf("failed to delete avatar in postgres: %w", err)
	}

	latestAvatar := s.repository.GetLatestSocietyAvatar(avatarInfo.UUID)

	err = s.produceNewSocietyAvatar(avatarInfo.UUID, latestAvatar)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to produce avatar: %v", err))
		return nil, fmt.Errorf("failed to produce avatar: %w", err)
	}

	return &avatarproto.Avatar{
		Id:   int32(avatarInfo.ID),
		Link: avatarInfo.Link,
	}, err
}

func (s *Service) uploadToS3(UUID, filename string, imageData []byte, avatarType AvatarType) (string, error) {
	objectName := fmt.Sprintf("%s/%s/%s", avatarType, UUID, generateTimestampedFileName(filename))
	contentType := "image/webp"

	link, err := s.s3Client.UploadFile(context.Background(), s.bucketName, objectName, imageData, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return link, nil
}

func generateTimestampedFileName(filename string) string {
	timestamp := time.Now().Format("20060102_150405")

	extension := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, extension)

	newExtension := ".webp"

	return fmt.Sprintf("%s_%s%s", timestamp, baseName, newExtension)
}
