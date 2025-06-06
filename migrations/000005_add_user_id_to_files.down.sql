-- Удаляем индексы
DROP INDEX IF EXISTS idx_encrypted_files_filename;
DROP INDEX IF EXISTS idx_encrypted_files_user_id;

-- Удаляем составной PRIMARY KEY
ALTER TABLE encrypted_files DROP CONSTRAINT encrypted_files_pkey;

-- Удаляем поле user_id
ALTER TABLE encrypted_files DROP COLUMN IF EXISTS user_id;

-- Восстанавливаем старый PRIMARY KEY на filename
ALTER TABLE encrypted_files ADD CONSTRAINT encrypted_files_pkey PRIMARY KEY (filename); 