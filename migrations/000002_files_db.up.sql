-- Создаем таблицу для хранения зашифрованных файлов
CREATE TABLE IF NOT EXISTS encrypted_files (
    filename TEXT PRIMARY KEY,
    file_data BYTEA NOT NULL,
    encrypted_key BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для поиска по имени файла
CREATE INDEX IF NOT EXISTS idx_encrypted_files_filename ON encrypted_files(filename);

-- Создаем триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_encrypted_files_updated_at
    BEFORE UPDATE ON encrypted_files
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 