package service

import (
	"context"
	"fmt"
	"io"

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

	err = s.updateAvatarInDB(userUUID, link)
	if err != nil {
		return err
	}

	err = s.produceAvatarRegistration(userUUID, link)
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
	objectName := fmt.Sprintf("%s/%s", userUUID, filename)
	contentType := "image/webp"

	link, err := s.s3Client.UploadFile(context.Background(), bucketName, objectName, imageData, contentType)
	if err != nil {
		return "", fmt.Errorf("error s.s3Client.UploadFile: %w", err)
	}

	return link, nil
}

func (s *Service) updateAvatarInDB(userUUID, link string) error {
	err := s.repository.SetAvatar(userUUID, link)
	if err != nil {
		return fmt.Errorf("error s.repository.SetAvatar: %w", err)
	}

	return nil
}

func (s *Service) produceAvatarRegistration(userUUID, link string) error {
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
