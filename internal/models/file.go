package models

import "time"

type FileMeta struct {
	Name      string
	UserID    uint // ID пользователя, которому принадлежит файл
	CreatedAt time.Time
	UpdatedAt time.Time
	// Поля для шифрования
	EncryptedKey []byte // Зашифрованный AES ключ
	Nonce        []byte // Вектор инициализации для AES-GCM
}

type FileReader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

// EncryptedFileData представляет зашифрованные данные файла
type EncryptedFileData struct {
	File     []byte // Зашифрованный файл
	Key      []byte // Зашифрованный AES ключ
	Nonce    []byte // Вектор инициализации
	Filename string // Оригинальное имя файла
}

