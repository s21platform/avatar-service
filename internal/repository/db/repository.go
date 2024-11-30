package db

import (
	"avatar_service/internal/config"
	"avatar_service/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL для использования в пакете database/sql
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) (*Repository, error) {
	var err error

	var repo *Repository

	for i := 0; i < 5; i++ {
		repo, err = connect(cfg)
		if err == nil {
			break
		}

		log.Println("failed to connect to database: ", err)
		time.Sleep(500 * time.Millisecond)
	}

	return repo, err
}

func connect(cfg *config.Config) (*Repository, error) {
	conStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
	)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &Repository{connection: db}, err
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) SetAvatar(userUUID, link string) error {
	query := `INSERT INTO avatar (user_uuid, link) VALUES ($1, $2)`
	_, err := r.connection.Exec(query, userUUID, link)

	if err != nil {
		return fmt.Errorf("failed to insert avatar into database: %w", err)
	}

	return nil
}

func (r *Repository) GetAllAvatars(userUUID string) (*model.AvatarInfoList, error) {
	var avatars model.AvatarInfoList

	query := `SELECT id, link FROM avatar WHERE user_uuid = $1 ORDER BY link DESC`

	err := r.connection.Select(&avatars, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch avatars from database: %w", err)
	}

	return &avatars, nil
}

func (r *Repository) GetAvatarData(avatarID int) (*model.AvatarInfo, error) {
	var avatarInfo model.AvatarInfo

	query := `SELECT id, user_uuid, link, create_at FROM avatar WHERE id = $1`
	err := r.connection.Get(&avatarInfo, query, avatarID)

	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	return &avatarInfo, nil
}

func (r *Repository) DeleteAvatar(avatarID int) error {
	query := `DELETE FROM avatar WHERE id = $1`
	_, err := r.connection.Exec(query, avatarID)

	if err != nil {
		return fmt.Errorf("failed to delete avatar from db: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestAvatar(userUUID string) string {
	var link string

	query := `SELECT link FROM avatar WHERE user_uuid = $1 ORDER BY create_at DESC LIMIT 1`

	err := r.connection.Get(&link, query, userUUID)
	if err != nil {
		return getDefaultAvatar()
	}

	return link
}

func getDefaultAvatar() string {
	return "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg"
}
