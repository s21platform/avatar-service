package postgres

import (
	"avatar_service/internal/config"
	"avatar_service/internal/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	return &Repository{
		connection: conn,
	}
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) SetUserAvatar(UUID, link string) error {
	query := `INSERT INTO users (uuid, link) VALUES ($1, $2)`
	_, err := r.connection.Exec(query, UUID, link)

	if err != nil {
		return fmt.Errorf("failed to insert avatar into database: %w", err)
	}

	return nil
}

func (r *Repository) GetAllUserAvatars(UUID string) (*model.AvatarInfoList, error) {
	var avatars model.AvatarInfoList

	query := `SELECT id, link FROM users WHERE uuid = $1 ORDER BY link DESC`

	err := r.connection.Select(&avatars, query, UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch avatars from database: %w", err)
	}

	return &avatars, nil
}

func (r *Repository) GetUserAvatarData(avatarID int) (*model.AvatarInfo, error) {
	var avatarInfo model.AvatarInfo

	query := `SELECT id, uuid, link, create_at FROM users WHERE id = $1`
	err := r.connection.Get(&avatarInfo, query, avatarID)

	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	return &avatarInfo, nil
}

func (r *Repository) DeleteUserAvatar(avatarID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.connection.Exec(query, avatarID)

	if err != nil {
		return fmt.Errorf("failed to delete avatar from postgres: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestUserAvatar(UUID string) string {
	var link string

	query := `SELECT link FROM users WHERE uuid = $1 ORDER BY create_at DESC LIMIT 1`

	err := r.connection.Get(&link, query, UUID)
	if err != nil {
		return getDefaultAvatar()
	}

	return link
}

func (r *Repository) SetSocietyAvatar(UUID, link string) error {
	query := `INSERT INTO society (uuid, link) VALUES ($1, $2)`
	_, err := r.connection.Exec(query, UUID, link)

	if err != nil {
		return fmt.Errorf("failed to insert avatar into database: %w", err)
	}

	return nil
}

func (r *Repository) GetAllSocietyAvatars(UUID string) (*model.AvatarInfoList, error) {
	var avatars model.AvatarInfoList

	query := `SELECT id, link FROM society WHERE uuid = $1 ORDER BY link DESC`

	err := r.connection.Select(&avatars, query, UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch avatars from database: %w", err)
	}

	return &avatars, nil
}

func (r *Repository) GetSocietyAvatarData(avatarID int) (*model.AvatarInfo, error) {
	var avatarInfo model.AvatarInfo

	query := `SELECT id, uuid, link, create_at FROM society WHERE id = $1`
	err := r.connection.Get(&avatarInfo, query, avatarID)

	if err != nil {
		return nil, fmt.Errorf("failed to get avatar data: %w", err)
	}

	return &avatarInfo, nil
}

func (r *Repository) DeleteSocietyAvatar(avatarID int) error {
	query := `DELETE FROM society WHERE id = $1`
	_, err := r.connection.Exec(query, avatarID)

	if err != nil {
		return fmt.Errorf("failed to delete avatar from postgres: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestSocietyAvatar(UUID string) string {
	var link string

	query := `SELECT link FROM society WHERE uuid = $1 ORDER BY create_at DESC LIMIT 1`

	err := r.connection.Get(&link, query, UUID)
	if err != nil {
		return getDefaultAvatar()
	}

	return link
}

func getDefaultAvatar() string {
	return "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg"
}
