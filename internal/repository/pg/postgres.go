package pg

import (
	"context"
	"fmt"

	"tages/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {

	return &Repository{pool: pool}
}

func (p *Repository) Close() {
	p.pool.Close()
}

func (p *Repository) SaveFileMeta(ctx context.Context, file *models.FileMeta) error {
	_, err := p.pool.Exec(ctx, SaveFileMetaQuery, file.Name, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save file meta for %s: %w", file.Name, err)
	}
	return nil
}

func (p *Repository) IsFileExists(ctx context.Context, filename string) (bool, error) {
	var exists bool
	err := p.pool.QueryRow(ctx, IsFileExistsQuery, filename).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of file %s: %w", filename, err)
	}
	return exists, nil
}

func (p *Repository) UpdateFileMeta(ctx context.Context, file *models.FileMeta) error {
	_, err := p.pool.Exec(ctx, UpdateFileMetaQuery, file.Name, file.CreatedAt, file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update file meta for %s: %w", file.Name, err)
	}
	return nil
}

func (p *Repository) GetFilesMeta(ctx context.Context) ([]*models.FileMeta, error) {
	rows, err := p.pool.Query(ctx, GetFilesMetaQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query files meta: %w", err)
	}
	defer rows.Close()

	var files []*models.FileMeta
	for rows.Next() {
		var file models.FileMeta
		if err := rows.Scan(&file.Name, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file meta row: %w", err)
		}
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return files, nil
}

func (p *Repository) DeleteFileMeta(ctx context.Context, filename string) error {

	_, err := p.pool.Exec(ctx, DeleteFileMetaQuery, filename)
	if err != nil {
		return fmt.Errorf("failed to delete file meta for %s: %w", filename, err)
	}
	return nil
}

// Методы для работы с пользователями
func (p *Repository) CreateUser(ctx context.Context, user *models.User) error {
	err := p.pool.QueryRow(ctx, CreateUserQuery,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (p *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := p.pool.QueryRow(ctx, GetUserByEmailQuery, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (p *Repository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := p.pool.QueryRow(ctx, GetUserByIDQuery, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}
