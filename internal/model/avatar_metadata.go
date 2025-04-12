package model

import (
	"time"

	avatarproto "github.com/s21platform/avatar-service/pkg/avatar"
)

type AvatarMetadataList []AvatarMetadata

type AvatarMetadata struct {
	ID        int       `db:"id"`
	UUID      string    `db:"uuid"`
	Link      string    `db:"link"`
	CreatedAt time.Time `db:"create_at"`
}

func (a *AvatarMetadataList) FromDTO() []*avatarproto.Avatar {
	result := make([]*avatarproto.Avatar, 0, len(*a))

	for _, avatar := range *a {
		result = append(result, &avatarproto.Avatar{
			Id:   int32(avatar.ID),
			Link: avatar.Link,
		})
	}

	return result
}
