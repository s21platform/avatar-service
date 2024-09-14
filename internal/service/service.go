package service

import (
	"context"
	"io"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type Service struct {
	avatarproto.UnimplementedAvatarServiceServer
	repository DBRepo
}

func New(repo DBRepo) *Service {
	return &Service{repository: repo}
}

func (s *Service) SetAvatar(stream avatarproto.AvatarService_SetAvatarServer) error {
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
			return err
		}

		if userUUID == "" && filename == "" {
			userUUID = in.UserUuid
			filename = in.Filename
		}

		imageData = append(imageData, in.Batch...)
	}

	link, err := s.repository.SetAvatar(userUUID, filename, imageData)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&avatarproto.SetAvatarOut{
		Link: link,
	})
}

func (s *Service) GetAllAvatars(
	ctx context.Context,
	in *avatarproto.GetAllAvatarsIn,
) (*avatarproto.GetAllAvatarsOut, error) {
	_ = ctx
	avatars, err := s.repository.GetAllAvatars(in.UserUuid)

	return &avatarproto.GetAllAvatarsOut{AvatarList: avatars}, err
}
