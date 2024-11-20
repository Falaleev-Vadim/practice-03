package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"strings"
)

func countLinesInFile(filename string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		resultChan <- fmt.Sprintf("Ошибка при открытии файла %s: %v", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		resultChan <- fmt.Sprintf("Ошибка при чтении файла %s: %v", filename, err)
		return
	}

	resultChan <- fmt.Sprintf("Файл %s содержит %d строк", filename, lineCount)
}

func findTextFiles(directory string) ([]string, error) {
	var textFiles []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			textFiles = append(textFiles, path)
		}
		return nil
	})
	return textFiles, err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <путь_к_каталогу>")
		return
	}

	directory := os.Args[1]

	textFiles, err := findTextFiles(directory)
	if err != nil {
		fmt.Printf("Ошибка при поиске файлов: %v\n", err)
		return
	}

	if len(textFiles) == 0 {
		fmt.Println("Не найдено файлов с расширением .txt")
		return
	}

	resultChan := make(chan string, len(textFiles))

	var wg sync.WaitGroup

	for _, file := range textFiles {
		wg.Add(1)
		go countLinesInFile(file, &wg, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println(result)
	}
}
