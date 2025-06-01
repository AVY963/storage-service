package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"

// 	"tages/internal/models"
// )

// type Repository struct {
// 	db *sql.DB
// }

// func New(db *sql.DB) *Repository {
// 	return &Repository{db: db}
// }

// func (r *Repository) UpdateFileMeta(ctx context.Context, fileMeta *models.FileMeta) error {
// 	query := `
// 		UPDATE encrypted_files
// 		SET updated_at = $2,
// 		    encrypted_key = $3,
// 		    nonce = $4
// 		WHERE filename = $1
// 	`

// 	result, err := r.db.ExecContext(ctx, query, fileMeta.Name, fileMeta.UpdatedAt, fileMeta.EncryptedKey, fileMeta.Nonce)
// 	if err != nil {
// 		return fmt.Errorf("failed to update file metadata: %w", err)
// 	}

// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %w", err)
// 	}

// 	if affected == 0 {
// 		return fmt.Errorf("file not found: %s", fileMeta.Name)
// 	}

// 	return nil
// }

// func (r *Repository) IsFileExists(ctx context.Context, filename string) (bool, error) {
// 	var exists bool
// 	query := `SELECT EXISTS(SELECT 1 FROM encrypted_files WHERE filename = $1)`

// 	err := r.db.QueryRowContext(ctx, query, filename).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check file existence: %w", err)
// 	}

// 	return exists, nil
// }

// func (r *Repository) SaveFileMeta(ctx context.Context, file *models.FileMeta) error {
// 	query := `
// 		INSERT INTO encrypted_files (filename, encrypted_key, nonce, created_at, updated_at)
// 		VALUES ($1, $2, $3, $4, $5)
// 		ON CONFLICT (filename) DO UPDATE
// 		SET encrypted_key = $2,
// 		    nonce = $3,
// 		    updated_at = $5
// 	`

// 	_, err := r.db.ExecContext(ctx, query, file.Name, file.EncryptedKey, file.Nonce, file.CreatedAt, file.UpdatedAt)
// 	if err != nil {
// 		return fmt.Errorf("failed to save file metadata: %w", err)
// 	}

// 	return nil
// }

// func (r *Repository) GetFilesMeta(ctx context.Context) ([]*models.FileMeta, error) {
// 	query := `
// 		SELECT filename, encrypted_key, nonce, created_at, updated_at
// 		FROM encrypted_files
// 		ORDER BY created_at DESC
// 	`

// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query files metadata: %w", err)
// 	}
// 	defer rows.Close()

// 	var files []*models.FileMeta
// 	for rows.Next() {
// 		var file models.FileMeta
// 		err := rows.Scan(
// 			&file.Name,
// 			&file.EncryptedKey,
// 			&file.Nonce,
// 			&file.CreatedAt,
// 			&file.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan file metadata: %w", err)
// 		}
// 		files = append(files, &file)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over rows: %w", err)
// 	}

// 	return files, nil
// }

// func (r *Repository) DeleteFileMeta(ctx context.Context, filename string) error {
// 	query := `DELETE FROM encrypted_files WHERE filename = $1`

// 	result, err := r.db.ExecContext(ctx, query, filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete file metadata: %w", err)
// 	}

// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %w", err)
// 	}

// 	if affected == 0 {
// 		return fmt.Errorf("file not found: %s", filename)
// 	}

// 	return nil
// }
