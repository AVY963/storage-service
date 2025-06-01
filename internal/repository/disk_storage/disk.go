package storage

// // вопрос насчет нейминга пакета
// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"tages/internal/models"
// 	"tages/internal/repository/pg"
// )

// type Storage struct {
// 	basePath string
// }

// func New(basePath string, db *pg.Repository) *Storage {

// 	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
// 		log.Printf("ERROR: Failed to create storage directory %s: %v", basePath, err)
// 	}

// 	return &Storage{
// 		basePath: basePath,
// 	}
// }

// func (ds *Storage) Save(filename string, data []byte) error {
// 	path := filepath.Join(ds.basePath, filename)
// 	log.Printf("INFO: Saving file to %s (%d bytes)", path, len(data))

// 	if err := os.WriteFile(path, data, 0644); err != nil {
// 		return fmt.Errorf("failed to save file to %s: %w", path, err)
// 	}

// 	log.Printf("INFO: File successfully saved to %s", path)
// 	return nil
// }

// func (ds *Storage) Read(filename string) ([]byte, error) {
// 	path := filepath.Join(ds.basePath, filename)
// 	log.Printf("INFO: Reading file from %s", path)

// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
// 	}

// 	log.Printf("INFO: Successfully read file %s (%d bytes)", path, len(data))
// 	return data, nil
// }

// func (ds *Storage) ReadStream(filename string) (models.FileReader, error) {
// 	path := filepath.Join(ds.basePath, filename)
// 	log.Printf("INFO: Opening file stream from %s", path)

// 	file, err := os.Open(path)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return nil, fmt.Errorf("file not found %s: %w", path, err)
// 		}
// 		return nil, fmt.Errorf("failed to open file stream for %s: %w", path, err)
// 	}

// 	log.Printf("INFO: Successfully opened file stream for %s", path)
// 	return file, nil
// }

// func (ds *Storage) Delete(filename string) error {
// 	path := filepath.Join(ds.basePath, filename)
// 	log.Printf("INFO: Deleting file: %s", path)

// 	err := os.Remove(path)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			log.Printf("WARNING: File not found for deletion: %s", path)
// 		} else {
// 			return fmt.Errorf("failed to delete fil %s: %w", path, err)
// 		}
// 	}

// 	log.Printf("INFO: File %s successfully deleted", path)
// 	return nil
// }
