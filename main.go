package main

import (
	"fmt"
	"os"
)

var DefaultRules = map[string]string{
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".pdf":  "Documents",
	".doc":  "Documents",
	".docx": "Documents",
	".txt":  "Documents",
	".mp3":  "Music",
	".wav":  "Music",
	".mp4":  "Video",
	".avi":  "Video",
	".zip":  "Archives",
	".rar":  "Archives",
}

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

// NewFileOrganizer создаёт новый FileOrganizer.
// Сигнатура: func NewFileOrganizer(sourceDir string) (*FileOrganizer, error)
func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if sourceDir == "" {
		return nil, fmt.Errorf("путь к директории не может быть пустым")
	}

	info, err := os.Stat(sourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("директория не существует: %s", sourceDir)
		}
		return nil, fmt.Errorf("не удалось получить информацию о директории %s: %w", sourceDir, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("указанный путь не является директорией: %s", sourceDir)
	}

	return &FileOrganizer{
		sourceDir: sourceDir,
		rulesMap:  DefaultRules,
	}, nil
}

func main() {
	for ext, category := range DefaultRules {
		fmt.Printf("%s -> %s\n", ext, category)
	}
}
