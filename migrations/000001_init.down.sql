-- Удаляем триггер
DROP TRIGGER IF EXISTS update_encrypted_files_updated_at ON encrypted_files;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем индекс
DROP INDEX IF EXISTS idx_encrypted_files_filename;

-- Удаляем таблицу
DROP TABLE IF EXISTS encrypted_files; 