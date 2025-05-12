package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"tages/internal/models"
)

type FileStorage interface {
	Save(filename string, data []byte) error
	Read(filename string) ([]byte, error)
	ReadStream(filename string) (models.FileReader, error)
	Delete(filename string) error
}

type Repository interface {
	UpdateFileMeta(ctx context.Context, filname *models.FileMeta) error
	IsFileExists(ctx context.Context, filename string) (bool, error)
	SaveFileMeta(ctx context.Context, file *models.FileMeta) error
	GetFilesMeta(ctx context.Context) ([]*models.FileMeta, error)
	DeleteFileMeta(ctx context.Context, filename string) error
}

type Usecase struct {
	storage FileStorage
	r       Repository
}

func New(storage FileStorage, r Repository) *Usecase {
	return &Usecase{
		storage: storage,
		r:       r,
	}
}
func (u *Usecase) Upload(ctx context.Context, filename string, data []byte) error {
	log.Printf("INFO: Processing upload request for file: %s (%d bytes)", filename, len(data))

	if err := u.storage.Save(filename, data); err != nil {
		return fmt.Errorf("failed to upload file %s: %w", filename, err)
	}

	exists, err := u.r.IsFileExists(ctx, filename)
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", filename, err)
	}

	now := time.Now()
	meta := &models.FileMeta{
		Name:      filename,
		UpdatedAt: now,
	}
	if !exists {
		meta.CreatedAt = now
	}

	if err := u.r.SaveFileMeta(ctx, meta); err != nil {
		return fmt.Errorf("failed to save file metadata for %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully uploaded file: %s", filename)
	return nil
}

func (u *Usecase) Download(filename string) (models.FileReader, error) {
	log.Printf("INFO: Processing download request for file: %s", filename)
	reader, err := u.storage.ReadStream(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to download file %s: %w", filename, err)
	}

	log.Printf("INFO: File stream opened for download: %s", filename)
	return reader, nil
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

	if err := u.storage.Delete(filename); err != nil {
		return fmt.Errorf("failed to delete file from storage %s: %w", filename, err)
	}

	if err := u.r.DeleteFileMeta(ctx, filename); err != nil {
		return fmt.Errorf("failed to delete file metadata for %s: %w", filename, err)
	}

	log.Printf("INFO: Successfully deleted file: %s", filename)
	return nil
}
