package service

import (
	"context"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type Service struct {
	avatarproto.UnimplementedAvatarServiceServer
	repository DBRepo
}

func New(repo DBRepo) *Service {
	return &Service{repository: repo}
}

//func (s *Service) SetAvatar(ctx context.Context, in *avatarproto.SetAvatarIn) (*avatarproto.SetAvatarOut, error) {
//	_ = ctx
//	link, err := s.repository.SetAvatar(in.UserUuid, in.Filename)
//
//	return &avatarproto.SetAvatarOut{Link: link}, err
//}

func (s *Service) GetAllAvatars(ctx context.Context, in *avatarproto.GetAllAvatarsIn) (*avatarproto.GetAllAvatarsOut, error) {
	_ = ctx
	avatars, err := s.repository.GetAllAvatars(in.UserUuid)

	return &avatarproto.GetAllAvatarsOut{AvatarList: avatars}, err
}
