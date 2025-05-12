package pg

import (
	"context"
	"fmt"
	"time"

	"tages/internal/models"
)

// Добавление refresh token в базу данных
func (p *Repository) StoreRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	_, err := p.pool.Exec(ctx, SaveRefreshTokenQuery,
		userID,
		token,
		expiresAt)

	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

// Удаление refresh token при выходе из системы
func (p *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := p.pool.Exec(ctx, DeleteRefreshTokenQuery, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

// Получение refresh token из базы данных
func (p *Repository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := p.pool.QueryRow(ctx, GetRefreshTokenQuery, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &refreshToken, nil
}
