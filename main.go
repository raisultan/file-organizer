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

func NewFileOrganizer(sourceDir string) *FileOrganizer {
	if sourceDir == "" {
		panic("путь к директории не может быть пустым")
	}

	info, err := os.Stat(sourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("директория не существует: %s", sourceDir))
		}
		panic(fmt.Sprintf("не удалось получить информацию о директории %s: %v", sourceDir, err))
	}

	if !info.IsDir() {
		panic(fmt.Sprintf("указанный путь не является директорией: %s", sourceDir))
	}

	return &FileOrganizer{
		sourceDir: sourceDir,
		rulesMap:  DefaultRules,
	}
}

func main() {
	for ext, category := range DefaultRules {
		fmt.Printf("%s -> %s\n", ext, category)
	}
}
