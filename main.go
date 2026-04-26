package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func (fo *FileOrganizer) initLog() error {
	logFile, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fo.logFile = logFile
	log.SetOutput(fo.logFile)
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	log.Printf("[SUCCESS] %s\n", message)
}

func (fo *FileOrganizer) logError(message string) {
	log.Printf("[ERROR] %s\n", message)
}

func (fo *FileOrganizer) Close() error {
	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	fullTargetDir := filepath.Join(fo.sourceDir, targetDir)
	if err := os.MkdirAll(fullTargetDir, 0755); err != nil {
		fo.logError(fmt.Sprintf("не удалось создать директорию %s: %v", fullTargetDir, err))
		return fmt.Errorf("не удалось создать директорию %s: %w", fullTargetDir, err)
	}

	fileName := filepath.Base(sourcePath)
	targetPath := filepath.Join(fullTargetDir, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		ext := filepath.Ext(fileName)
		name := strings.TrimSuffix(fileName, ext)
		timestamp := time.Now().Format("_2006-01-02_15-04-05")
		fileName = name + timestamp + ext
		targetPath = filepath.Join(fullTargetDir, fileName)
	}

	if err := os.Rename(sourcePath, targetPath); err != nil {
		fo.logError(fmt.Sprintf("не удалось переместить файл %s в %s: %v", sourcePath, targetPath, err))
		return fmt.Errorf("не удалось переместить файл %s в %s: %w", sourcePath, targetPath, err)
	}

	fo.logSuccess(fmt.Sprintf("перемещён %s -> %s", sourcePath, targetPath))
	return nil
}

func main() {
	for ext, category := range DefaultRules {
		fmt.Printf("%s -> %s\n", ext, category)
	}
}
