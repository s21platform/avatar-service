package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
	"github.com/s21platform/user-proto/user-proto/new_avatar_register"
)

type Service struct {
	avatarproto.UnimplementedAvatarServiceServer
	s3Client      S3Storage
	repository    DBRepo
	kafkaProducer NewAvatarRegisterSrv
}

func New(s3Client S3Storage, repo DBRepo, kafkaProducer NewAvatarRegisterSrv) *Service {
	return &Service{
		s3Client:      s3Client,
		repository:    repo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *Service) SetAvatar(stream avatarproto.AvatarService_SetAvatarServer) error {
	userUUID, filename, imageData, err := s.receiveData(stream)
	if err != nil {
		return err
	}

	link, err := s.uploadToS3(userUUID, filename, imageData)
	if err != nil {
		return err
	}

	err = s.uploadToDB(userUUID, link)
	if err != nil {
		return err
	}

	err = s.produceNewAvatar(userUUID, link)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&avatarproto.SetAvatarOut{
		Link: link,
	})
}

func (s *Service) receiveData(stream avatarproto.AvatarService_SetAvatarServer) (string, string, []byte, error) {
	var (
		userUUID  string
		filename  string
		imageData []byte
	)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", "", nil, fmt.Errorf("error stream.Recv: %w", err)
		}

		if userUUID == "" && filename == "" {
			userUUID = in.UserUuid
			filename = in.Filename
		}

		imageData = append(imageData, in.Batch...)
	}

	return userUUID, filename, imageData, nil
}

func (s *Service) uploadToS3(userUUID, filename string, imageData []byte) (string, error) {
	bucketName := "space21staging"
	objectName := getObjectName(userUUID, filename)
	contentType := "image/webp"

	link, err := s.s3Client.UploadFile(context.Background(), bucketName, objectName, imageData, contentType)
	if err != nil {
		return "", fmt.Errorf("error s.s3Client.UploadFile: %w", err)
	}

	return link, nil
}

func getObjectName(userUUID, filename string) string {
	return fmt.Sprintf("%s/%s", userUUID, generateTimestampedFileName(filename))
}

func generateTimestampedFileName(filename string) string {
	timestamp := time.Now().Format("20060102_150405")

	extension := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, extension)

	newExtension := ".webp"

	return fmt.Sprintf("%s_%s%s", timestamp, baseName, newExtension)
}

func (s *Service) uploadToDB(userUUID, link string) error {
	err := s.repository.SetAvatar(userUUID, link)
	if err != nil {
		return fmt.Errorf("error s.repository.SetAvatar: %w", err)
	}

	return nil
}

func (s *Service) produceNewAvatar(userUUID, link string) error {
	msg := &new_avatar_register.NewAvatarRegister{
		Uuid: userUUID,
		Link: link,
	}

	err := s.kafkaProducer.ProduceMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAllAvatars(ctx context.Context,
	in *avatarproto.GetAllAvatarsIn) (*avatarproto.GetAllAvatarsOut, error) {
	_ = ctx
	avatars, err := s.repository.GetAllAvatars(in.UserUuid)

	return &avatarproto.GetAllAvatarsOut{AvatarList: avatars}, err
}

func (s *Service) DeleteAvatar(ctx context.Context, in *avatarproto.DeleteAvatarIn) (*avatarproto.Avatar, error) {
	avatarInfo, err := s.repository.GetAvatarData(int(in.AvatarId))
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	err = s.s3Client.DeleteAvatar(ctx, avatarInfo.Link)
	if err != nil {
		return nil, fmt.Errorf("failed to delete avatar in s3: %w", err)
	}

	err = s.repository.DeleteAvatar(avatarInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete avatar in db: %w", err)
	}

	latestAvatar := s.repository.GetLatestAvatar(avatarInfo.UserUUID)
	err = s.produceNewAvatar(avatarInfo.UserUUID, latestAvatar)

	if err != nil {
		return nil, fmt.Errorf("failed to produce avatar: %w", err)
	}

	return &avatarproto.Avatar{
		//nolint: gosec
		Id:   int32(avatarInfo.ID),
		Link: avatarInfo.Link,
	}, err
}
