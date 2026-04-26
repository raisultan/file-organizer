package main

import (
	"bufio"
	"fmt"
	"io/fs"
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

type FileStats struct {
	Count     int
	TotalSize int64
}

func (fs *FileStats) String() string {
	return fmt.Sprintf("Файлов: %d, Размер: %.2f KB", fs.Count, float64(fs.TotalSize)/1024)
}

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
	statistics     map[string]*FileStats
	totalSize      int64
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
		sourceDir:  sourceDir,
		rulesMap:   DefaultRules,
		statistics: make(map[string]*FileStats),
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

func (fo *FileOrganizer) Organize() error {
	if err := fo.initLog(); err != nil {
		return fmt.Errorf("не удалось инициализировать лог: %w", err)
	}

	return filepath.WalkDir(fo.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fo.logError(fmt.Sprintf("ошибка обхода %s: %v", path, err))
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Dir(path) != fo.sourceDir {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		targetDir, ok := fo.rulesMap[ext]
		if !ok {
			return nil
		}

		info, infoErr := d.Info()
		if infoErr != nil {
			fo.logError(fmt.Sprintf("не удалось получить информацию о файле %s: %v", path, infoErr))
			return nil
		}

		if err := fo.moveFile(path, targetDir); err != nil {
			return nil
		}

		if _, exists := fo.statistics[targetDir]; !exists {
			fo.statistics[targetDir] = &FileStats{}
		}
		fo.statistics[targetDir].Count++
		fo.statistics[targetDir].TotalSize += info.Size()
		fo.totalSize += info.Size()
		fo.processedFiles++
		return nil
	})
}

func (fo *FileOrganizer) generateReport() string {
	var b strings.Builder
	b.WriteString("=== Отчёт о перемещении файлов ===\n\n")
	b.WriteString(fmt.Sprintf("Всего обработано файлов: %d\n", fo.processedFiles))
	b.WriteString(fmt.Sprintf("Общий размер: %.2f KB\n\n", float64(fo.totalSize)/1024))
	b.WriteString("Статистика по категориям:\n\n")
	for category, stats := range fo.statistics {
		b.WriteString(fmt.Sprintf("%s:\n  %s\n\n", category, stats))
	}
	return b.String()
}

func main() {
	fmt.Println("=== Файловый органайзер ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите путь к директории для организации (Enter для текущей директории): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Ошибка чтения ввода: %v\n", err)
		os.Exit(1)
	}

	sourcePath := strings.TrimSpace(input)
	if sourcePath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Не удалось определить текущую директорию: %v\n", err)
			os.Exit(1)
		}
		sourcePath = cwd
	}

	organizer, err := NewFileOrganizer(sourcePath)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
	defer organizer.Close()

	if err := organizer.Organize(); err != nil {
		fmt.Printf("Ошибка при организации файлов: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Print(organizer.generateReport())
	fmt.Println("Лог операций: organizer.log")
}
