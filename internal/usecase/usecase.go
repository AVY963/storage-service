package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"tages/internal/models"
)

type Repository interface {
	SaveFile(filename string, data []byte, key []byte, nonce []byte, userID uint) error
	ReadFile(filename string, userID uint) (*models.EncryptedFileData, error)
	DeleteFile(ctx context.Context, filename string, userID uint) error
	UpdateFileMeta(ctx context.Context, fileMeta *models.FileMeta) error
	IsFileExists(ctx context.Context, filename string, userID uint) (bool, error)
	SaveFileMeta(ctx context.Context, file *models.FileMeta) error
	GetFilesMeta(ctx context.Context, userID uint) ([]*models.FileMeta, error)
	DeleteFileMeta(ctx context.Context, filename string, userID uint) error
}

type Usecase struct {
	r Repository
}

func New(r Repository) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (u *Usecase) Upload(ctx context.Context, data *models.EncryptedFileData, userID uint) error {
	log.Printf("INFO: Processing upload request for file: %s (user: %d)", data.Filename, userID)

	if err := u.r.SaveFile(data.Filename, data.File, data.Key, data.Nonce, userID); err != nil {
		return fmt.Errorf("failed to upload file %s: %w", data.Filename, err)
	}

	exists, err := u.r.IsFileExists(ctx, data.Filename, userID)
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", data.Filename, err)
	}

	now := time.Now()
	meta := &models.FileMeta{
		Name:         data.Filename,
		UserID:       userID,
		UpdatedAt:    now,
		EncryptedKey: data.Key,
		Nonce:        data.Nonce,
	}
	if !exists {
		meta.CreatedAt = now
	}

	if err := u.r.SaveFileMeta(ctx, meta); err != nil {
		return fmt.Errorf("failed to save file metadata for %s: %w", data.Filename, err)
	}

	log.Printf("INFO: Successfully uploaded file: %s (user: %d)", data.Filename, userID)
	return nil
}

func (u *Usecase) Download(ctx context.Context, filename string, userID uint) (*models.EncryptedFileData, error) {
	log.Printf("INFO: Processing download request for file: %s (user: %d)", filename, userID)

	data, err := u.r.ReadFile(filename, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to download file %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully downloaded file: %s (user: %d)", filename, userID)
	return data, nil
}

func (u *Usecase) ListFiles(ctx context.Context, userID uint) ([]*models.FileMeta, error) {
	log.Printf("INFO: Retrieving file list for user: %d", userID)

	files, err := u.r.GetFilesMeta(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve file list for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve file list: %w", err)
	}

	log.Printf("INFO: Retrieved %d files for user %d", len(files), userID)
	return files, nil
}

func (u *Usecase) DeleteFile(ctx context.Context, filename string, userID uint) error {
	log.Printf("INFO: Processing delete request for file: %s (user: %d)", filename, userID)

	if err := u.r.DeleteFile(ctx, filename, userID); err != nil {
		return fmt.Errorf("failed to delete file from storage %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully deleted file: %s (user: %d)", filename, userID)
	return nil
}
