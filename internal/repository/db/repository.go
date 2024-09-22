package db

import (
	"avatar_service/internal/config"
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

func (r *Repository) GetAllAvatars(userUUID string) ([]string, error) {
	var avatars []string

	err := r.connection.Select(&avatars, `SELECT link FROM avatar WHERE user_uuid = $1`, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error r.connection.Select: %w", err)
	}

	return avatars, nil
}
