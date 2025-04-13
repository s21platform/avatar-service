package postgres

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL

	"github.com/s21platform/avatar-service/internal/config"
	"github.com/s21platform/avatar-service/internal/model"
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

func (r *Repository) SetUserAvatar(ctx context.Context, uuid, link string) error {
	query, args, err := sq.
		Insert("users").
		Columns("uuid", "link").
		Values(uuid, link).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert user avatar: %w", err)
	}

	return nil
}

func (r *Repository) GetAllUserAvatars(ctx context.Context, uuid string) (*model.AvatarMetadataList, error) {
	var avatars model.AvatarMetadataList

	query, args, err := sq.
		Select("id", "link").
		From("users").
		Where(sq.Eq{"uuid": uuid}).
		OrderBy("link DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.connection.SelectContext(ctx, &avatars, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user avatars: %w", err)
	}

	return &avatars, nil
}

func (r *Repository) GetUserAvatarData(ctx context.Context, avatarID int) (*model.AvatarMetadata, error) {
	var avatarInfo model.AvatarMetadata

	query, args, err := sq.
		Select("id", "uuid", "link", "create_at").
		From("users").
		Where(sq.Eq{"id": avatarID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.connection.GetContext(ctx, &avatarInfo, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user avatar data: %w", err)
	}

	return &avatarInfo, nil
}

func (r *Repository) DeleteUserAvatar(ctx context.Context, avatarID int) error {
	query, args, err := sq.
		Delete("users").
		Where(sq.Eq{"id": avatarID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user avatar: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestUserAvatar(ctx context.Context, uuid string) string {
	var link string

	query, args, err := sq.
		Select("link").
		From("users").
		Where(sq.Eq{"uuid": uuid}).
		OrderBy("create_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return getDefaultAvatar()
	}

	err = r.connection.GetContext(ctx, &link, query, args...)
	if err != nil {
		return getDefaultAvatar()
	}

	return link
}

func (r *Repository) SetSocietyAvatar(ctx context.Context, uuid, link string) error {
	query, args, err := sq.
		Insert("society").
		Columns("uuid", "link").
		Values(uuid, link).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert society avatar: %w", err)
	}

	return nil
}

func (r *Repository) GetAllSocietyAvatars(ctx context.Context, uuid string) (*model.AvatarMetadataList, error) {
	var avatars model.AvatarMetadataList

	query, args, err := sq.
		Select("id", "link").
		From("society").
		Where(sq.Eq{"uuid": uuid}).
		OrderBy("link DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.connection.SelectContext(ctx, &avatars, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch society avatars: %w", err)
	}

	return &avatars, nil
}

func (r *Repository) GetSocietyAvatarData(ctx context.Context, avatarID int) (*model.AvatarMetadata, error) {
	var avatarInfo model.AvatarMetadata

	query, args, err := sq.
		Select("id", "uuid", "link", "create_at").
		From("society").
		Where(sq.Eq{"id": avatarID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	err = r.connection.GetContext(ctx, &avatarInfo, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get society avatar data: %w", err)
	}

	return &avatarInfo, nil
}

func (r *Repository) DeleteSocietyAvatar(ctx context.Context, avatarID int) error {
	query, args, err := sq.
		Delete("society").
		Where(sq.Eq{"id": avatarID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete society avatar: %w", err)
	}

	return nil
}

func (r *Repository) GetLatestSocietyAvatar(ctx context.Context, uuid string) string {
	var link string

	query, args, err := sq.
		Select("link").
		From("society").
		Where(sq.Eq{"uuid": uuid}).
		OrderBy("create_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return getDefaultAvatar()
	}

	err = r.connection.GetContext(ctx, &link, query, args...)
	if err != nil {
		return getDefaultAvatar()
	}

	return link
}

func getDefaultAvatar() string {
	return "https://storage.yandexcloud.net/space21/avatars/default/logo-discord.jpeg"
}
