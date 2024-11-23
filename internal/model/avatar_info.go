package model

import (
	"time"

	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
)

type AvatarInfoList []AvatarInfo

type AvatarInfo struct {
	ID        int       `db:"id"`
	UserUUID  string    `db:"user_uuid"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"create_at"`
}

func (a *AvatarInfoList) FromDTO() []*avatarproto.Avatar {
	result := make([]*avatarproto.Avatar, 0, len(*a))

	for _, avatar := range *a {
		result = append(result, avatar.ToProto())
	}

	return result
}

func (a *AvatarInfo) ToProto() *avatarproto.Avatar {
	return &avatarproto.Avatar{
		//nolint: gosec
		Id:   int32(a.ID),
		Link: a.Link,
	}
}
