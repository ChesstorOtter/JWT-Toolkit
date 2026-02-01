package jwt

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/schollz/progressbar/v3"
)

func CrackHS256(tokenString string, worldlistPath string) {
	file, err := os.Open(worldlistPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening wordlist file: %v\n", err)
		return
	}
	defer file.Close()

	lineCount := countLines(worldlistPath)
	bar := progressbar.Default(int64(lineCount))

	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	secrets := make(chan string, numWorkers*2)
	found := atomic.Bool{}
	foundSecret := ""
	progressCount := int64(0)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for secret := range secrets {
				atomic.AddInt64(&progressCount, 1)
				bar.Add(1)

				if found.Load() {
					continue
				}

				valid, err := VerifyHS256(tokenString, []byte(secret))
				if err != nil {
					continue
				}

				if valid {
					found.Store(true)
					foundSecret = secret
					return
				}
			}
		}()
	}

	go func() {
		file.Seek(0, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if found.Load() {
				break
			}
			secret := scanner.Text()
			secrets <- secret
		}
		close(secrets)
	}()

	wg.Wait()
	bar.Finish()

	if found.Load() {
		fmt.Println("\nToken is valid with secret:", foundSecret)
	} else {
		fmt.Println("\nNo valid secret found")
	}
}

func countLines(filePath string) int {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
