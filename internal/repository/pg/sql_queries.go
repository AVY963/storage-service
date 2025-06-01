package pg

const (
	// Запросы для работы с файлами
	SaveFileQuery = `
		INSERT INTO encrypted_files (filename, file_data, encrypted_key, nonce, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (filename) DO UPDATE
		SET file_data = $2, 
		    encrypted_key = $3, 
		    nonce = $4,
		    updated_at = $6
	`

	ReadFileQuery = `
		SELECT file_data, encrypted_key, nonce
		FROM encrypted_files
		WHERE filename = $1
	`

	DeleteFileQuery = `
		DELETE FROM encrypted_files WHERE filename = $1
	`

	// Запросы для работы с метаданными файлов
	IsFileExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM encrypted_files WHERE filename = $1)
	`

	GetFilesMetaQuery = `
		SELECT filename, encrypted_key, nonce, created_at, updated_at 
		FROM encrypted_files 
		ORDER BY created_at DESC
	`

	UpdateFileMetaQuery = `
		UPDATE encrypted_files
		SET encrypted_key = $2,
		    nonce = $3,
		    updated_at = $4
		WHERE filename = $1
	`

	DeleteFileMetaQuery = `
		DELETE FROM encrypted_files WHERE filename = $1
	`

	// Запросы для пользователей
	CreateUserQuery = `
		INSERT INTO users(email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	GetUserByEmailQuery = `
		SELECT id, email, password, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	GetUserByIDQuery = `
		SELECT id, email, password, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	UpdateUserQuery = `
		UPDATE users
		SET updated_at = $2
		WHERE id = $1
	`

	// Запросы для refresh токенов
	SaveRefreshTokenQuery = `
		INSERT INTO refresh_tokens(user_id, token, expires_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE 
		SET token = $2, expires_at = $3
	`

	GetRefreshTokenQuery = `
		SELECT id, user_id, token, expires_at 
		FROM refresh_tokens 
		WHERE token = $1 AND expires_at > NOW()
	`

	DeleteRefreshTokenQuery = `
		DELETE FROM refresh_tokens WHERE token = $1
	`
)
