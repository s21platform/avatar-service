package model

type AvatarType string

const (
	UserAvatarType    AvatarType = "user"
	SocietyAvatarType AvatarType = "society"
)

type AvatarContent struct {
	AvatarType AvatarType
	UUID       string
	Filename   string
	ImageData  []byte
}
