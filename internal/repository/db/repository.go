package db

import (
	"avatar_service/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL для использования в пакете database/sql
)

type Repository struct {
	connection *sql.DB
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

	db, err := sql.Open("postgres", conStr)

	if err != nil {
		return nil, fmt.Errorf("error sql.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error db.Ping: %w", err)
	}

	return &Repository{db}, err
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}
