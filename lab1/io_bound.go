package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// I/O-BOUND

func countWordsInFile(path string) int64 {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return int64(len(strings.Fields(string(data))))
}

// runIOBound recursively iterate through the directory and count words
// limited by speed of hard drive
func runIOBound(root string, parallel bool) {
	start := time.Now()
	var totalWords int64

	var files []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if !parallel {
		for _, f := range files {
			totalWords += countWordsInFile(f)
		}

		fmt.Printf("I/O-Bound (sequential): words: %d, time: %v\n",
			totalWords, time.Since(start))

		return
	}

	var wg sync.WaitGroup
	wordChan := make(chan int64, len(files))
	semaphore := make(chan struct{}, 20)

	for _, f := range files {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			semaphore <- struct{}{}
			wordChan <- countWordsInFile(p)
			<-semaphore
		}(f)
	}
	go func() {
		wg.Wait()
		close(wordChan)
	}()

	for n := range wordChan {
		totalWords += n
	}

	fmt.Printf("I/O-Bound (parallel): words: %d, time: %v\n",
		totalWords, time.Since(start))
}

// generate files
func generateRandomWord(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func createRandomFiles(root string, count int64) {
	rand.Seed(time.Now().UnixNano())

	for i := int64(0); i < count; i++ {
		subDir := filepath.Join(root, fmt.Sprintf("dir_%d", rand.Intn(5)))
		_ = os.MkdirAll(subDir, 0755)

		fileName := fmt.Sprintf("file_%d.txt", i)
		filePath := filepath.Join(subDir, fileName)

		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Помилка створення файлу %s: %v\n", filePath, err)
			continue
		}

		wordCount := rand.Intn(90) + 10
		for j := 0; j < wordCount; j++ {
			_, _ = f.WriteString(generateRandomWord(5) + " ")
		}
		f.Close()
	}
}
