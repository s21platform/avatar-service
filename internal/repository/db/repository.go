package db

import (
	"avatar_service/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL для использования в пакете database/sql
	avatarproto "github.com/s21platform/avatar-proto/avatar-proto"
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

		log.Println("error connect(cfg) ", err)
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
		return nil, fmt.Errorf("error sql.Open: %w", err)
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
		return fmt.Errorf("error r.connection.Exec: %w", err)
	}

	return nil
}

func (r *Repository) GetAllAvatars(userUUID string) ([]*avatarproto.Avatar, error) {
	var avatars []*avatarproto.Avatar

	query := `SELECT id, link FROM avatar WHERE user_uuid = $1 ORDER BY link DESC`

	err := r.connection.Select(&avatars, query, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error r.connection.Select: %w", err)
	}

	return avatars, nil
}

func (r *Repository) GetAvatarData(avatarID int) (int, string, string, time.Time, error) {
	var avatarData struct {
		ID        int       `db:"id"`
		UserUUID  string    `db:"user_uuid"`
		Link      string    `db:"link"`
		CreatedAt time.Time `db:"create_at"`
	}

	query := `SELECT id, user_uuid, link, create_at FROM avatar WHERE id = $1`
	err := r.connection.Get(&avatarData, query, avatarID)

	if err != nil {
		return 0, "", "", time.Time{}, fmt.Errorf("error r.connection.Get: %w", err)
	}

	return avatarData.ID, avatarData.UserUUID, avatarData.Link, avatarData.CreatedAt, nil
}

func (r *Repository) DeleteAvatar(avatarID int) error {
	query := `DELETE FROM avatar WHERE id = $1`
	_, err := r.connection.Exec(query, avatarID)

	if err != nil {
		return fmt.Errorf("error r.connection.Exec: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestAvatar(userUUID string) string {
	var link string

	query := `SELECT link FROM avatar WHERE user_uuid = $1 ORDER BY create_at DESC LIMIT 1`

	err := r.connection.Get(&link, query, userUUID)
	if err != nil {
		return ""
	}

	return link
}
