package service

type DBRepo interface {
	//SetAvatar(userUuid, filename string) (string, error)
	GetAllAvatars(userUuid string) ([]string, error)
}
