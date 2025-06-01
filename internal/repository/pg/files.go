package pg

import (
	"context"
	"fmt"
	"time"

	"tages/internal/models"
)

// SaveFile сохраняет зашифрованный файл в базу данных
func (p *Repository) SaveFile(filename string, data []byte, key []byte, nonce []byte) error {
	now := time.Now()
	_, err := p.pool.Exec(context.Background(), SaveFileQuery,
		filename,
		data,
		key,
		nonce,
		now, // created_at
		now, // updated_at
	)
	if err != nil {
		return fmt.Errorf("failed to save file to database: %w", err)
	}
	return nil
}

// ReadFile читает зашифрованный файл из базы данных
func (p *Repository) ReadFile(filename string) (*models.EncryptedFileData, error) {
	var data models.EncryptedFileData
	data.Filename = filename

	err := p.pool.QueryRow(context.Background(), ReadFileQuery, filename).Scan(&data.File, &data.Key, &data.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from database: %w", err)
	}

	return &data, nil
}

// DeleteFile удаляет файл из базы данных
func (p *Repository) DeleteFile(ctx context.Context, filename string) error {
	result, err := p.pool.Exec(ctx, DeleteFileQuery, filename)
	if err != nil {
		return fmt.Errorf("failed to delete file meta for %s: %w", filename, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found: %s", filename)
	}

	return nil
}


// Методы для работы с метаданными

func (p *Repository) SaveFileMeta(ctx context.Context, file *models.FileMeta) error {
	exists, err := p.IsFileExists(ctx, file.Name)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("file data must be saved before metadata: %s", file.Name)
	}

	result, err := p.pool.Exec(ctx, UpdateFileMetaQuery,
		file.Name,
		file.EncryptedKey,
		file.Nonce,
		file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save file meta for %s: %w", file.Name, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found: %s", file.Name)
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
	result, err := p.pool.Exec(ctx, UpdateFileMetaQuery,
		file.Name,
		file.EncryptedKey,
		file.Nonce,
		file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update file meta for %s: %w", file.Name, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found: %s", file.Name)
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
		err := rows.Scan(
			&file.Name,
			&file.EncryptedKey,
			&file.Nonce,
			&file.CreatedAt,
			&file.UpdatedAt,
		)
		if err != nil {
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
	result, err := p.pool.Exec(ctx, DeleteFileMetaQuery, filename)
	if err != nil {
		return fmt.Errorf("failed to delete file meta for %s: %w", filename, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("file not found: %s", filename)
	}

	return nil
}
