package pg

import (
	"context"
	"fmt"

	"tages/internal/models"
)

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

// UpdateUser обновляет информацию о пользователе
func (p *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := p.pool.Exec(ctx, UpdateUserQuery,
		user.ID,
		user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
