package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"tages/internal/models"
)

type Repository interface {
	SaveFile(filename string, data []byte, key []byte, nonce []byte) error
	ReadFile(filename string) (*models.EncryptedFileData, error)
	DeleteFile(ctx context.Context, filename string) error
	UpdateFileMeta(ctx context.Context, fileMeta *models.FileMeta) error
	IsFileExists(ctx context.Context, filename string) (bool, error)
	SaveFileMeta(ctx context.Context, file *models.FileMeta) error
	GetFilesMeta(ctx context.Context) ([]*models.FileMeta, error)
	DeleteFileMeta(ctx context.Context, filename string) error
}

type Usecase struct {
	r Repository
}

func New(r Repository) *Usecase {
	return &Usecase{
		r: r,
	}
}

func (u *Usecase) Upload(ctx context.Context, data *models.EncryptedFileData) error {
	log.Printf("INFO: Processing upload request for file: %s", data.Filename)

	if err := u.r.SaveFile(data.Filename, data.File, data.Key, data.Nonce); err != nil {
		return fmt.Errorf("failed to upload file %s: %w", data.Filename, err)
	}

	exists, err := u.r.IsFileExists(ctx, data.Filename)
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", data.Filename, err)
	}

	now := time.Now()
	meta := &models.FileMeta{
		Name:         data.Filename,
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

	log.Printf("INFO: Successfully uploaded file: %s", data.Filename)
	return nil
}

func (u *Usecase) Download(ctx context.Context, filename string) (*models.EncryptedFileData, error) {
	log.Printf("INFO: Processing download request for file: %s", filename)

	data, err := u.r.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to download file %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully downloaded file: %s", filename)
	log.Printf("INFO: File data: %v", data)
	return data, nil
}

func (u *Usecase) ListFiles(ctx context.Context) ([]*models.FileMeta, error) {
	log.Printf("INFO: Retrieving file list")

	files, err := u.r.GetFilesMeta(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve file list: %v", err)
		return nil, fmt.Errorf("failed to retrieve file list: %w", err)
	}

	log.Printf("INFO: Retrieved %d files", len(files))
	return files, nil
}

func (u *Usecase) DeleteFile(ctx context.Context, filename string) error {
	log.Printf("INFO: Processing delete request for file: %s", filename)

	if err := u.r.DeleteFile(ctx, filename); err != nil {
		return fmt.Errorf("failed to delete file from storage %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully deleted file: %s", filename)
	return nil
}
