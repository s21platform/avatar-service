package service

type DBRepo interface {
	SetAvatar(userUUID, filename string, imageData []byte) (string, error)
	GetAllAvatars(userUUID string) ([]string, error)
}
