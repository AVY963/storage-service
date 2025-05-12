package pg

const (
	SaveFileMetaQuery = `
		INSERT INTO file_meta(name, created_at, updated_at) 
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO UPDATE 
		SET updated_at = $3
	`
	IsFileExistsQuery = `
		SELECT EXISTS(SELECT 1 FROM file_meta WHERE name = $1)
	`

	GetFilesMetaQuery = `
		SELECT name, created_at, updated_at 
		FROM file_meta 
		ORDER BY updated_at DESC
	`

	UpdateFileMetaQuery = `
		UPDATE file_meta
		SET created_at = $2, updated_at = $3
		WHERE name = $1
	`

	DeleteFileMetaQuery = `
		DELETE FROM file_meta WHERE name = $1
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
