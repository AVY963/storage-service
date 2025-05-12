package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"tages/internal/auth"
	"tages/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	StoreRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, token string) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
}

type AuthManager interface {
	GenerateTokenPair(user *models.User) (*auth.TokenPair, error)
	ParseAccessToken(tokenString string) (*auth.TokenClaims, error)
	ParseRefreshToken(tokenString string) (*auth.TokenClaims, error)
}

type UserUsecase struct {
	userRepo    UserRepository
	authManager AuthManager
}

// Ошибки авторизации
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid refresh token")
)

func NewUserUsecase(userRepo UserRepository, authManager AuthManager) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		authManager: authManager,
	}
}

// Регистрация нового пользователя
func (u *UserUsecase) Register(ctx context.Context, email, password string) (*models.User, *auth.TokenPair, error) {
	log.Printf("INFO: Registering new user with email: %s", email)

	// Проверяем, существует ли пользователь с таким email
	_, err := u.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, nil, ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем нового пользователя
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем токены
	tokens, err := u.authManager.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Сохраняем refresh token в базу
	claims, err := u.authManager.ParseRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	if err := u.userRepo.StoreRefreshToken(ctx, user.ID, tokens.RefreshToken, expiresAt); err != nil {
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	log.Printf("INFO: Successfully registered user with ID: %d", user.ID)
	return user, tokens, nil
}

// Авторизация пользователя
func (u *UserUsecase) Login(ctx context.Context, email, password string) (*models.User, *auth.TokenPair, error) {
	log.Printf("INFO: Login attempt for user: %s", email)

	// Получаем пользователя по email
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Генерируем токены
	tokens, err := u.authManager.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Сохраняем refresh token в базу
	claims, err := u.authManager.ParseRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	if err := u.userRepo.StoreRefreshToken(ctx, user.ID, tokens.RefreshToken, expiresAt); err != nil {
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	log.Printf("INFO: User %s logged in successfully", email)
	return user, tokens, nil
}

// Обновление токена
func (u *UserUsecase) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	log.Printf("INFO: Refreshing token")

	// Проверяем refresh token
	claims, err := u.authManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Получаем сохраненный токен из БД
	storedToken, err := u.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Получаем пользователя
	user, err := u.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Генерируем новую пару токенов
	tokens, err := u.authManager.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Удаляем старый refresh token
	if err := u.userRepo.DeleteRefreshToken(ctx, storedToken.Token); err != nil {
		log.Printf("WARN: Failed to delete old refresh token: %v", err)
	}

	// Сохраняем новый refresh token
	claims, err = u.authManager.ParseRefreshToken(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	if err := u.userRepo.StoreRefreshToken(ctx, user.ID, tokens.RefreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	log.Printf("INFO: Successfully refreshed token for user ID: %d", user.ID)
	return tokens, nil
}

// Выход из системы
func (u *UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	log.Printf("INFO: Processing logout")

	if err := u.userRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	log.Printf("INFO: User successfully logged out")
	return nil
}

// Получение пользователя по access token
func (u *UserUsecase) GetUserByToken(ctx context.Context, accessToken string) (*models.User, error) {
	claims, err := u.authManager.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
