-- Добавляем поле user_id в таблицу encrypted_files
ALTER TABLE encrypted_files 
ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;

-- Удаляем старый PRIMARY KEY
ALTER TABLE encrypted_files DROP CONSTRAINT encrypted_files_pkey;

-- Создаем новый составной PRIMARY KEY
ALTER TABLE encrypted_files ADD CONSTRAINT encrypted_files_pkey PRIMARY KEY (filename, user_id);

-- Создаем индекс для быстрого поиска файлов пользователя
CREATE INDEX IF NOT EXISTS idx_encrypted_files_user_id ON encrypted_files(user_id);

-- Создаем индекс для поиска по filename (для совместимости)
CREATE INDEX IF NOT EXISTS idx_encrypted_files_filename ON encrypted_files(filename); 