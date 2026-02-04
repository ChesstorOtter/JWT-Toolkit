package jwt

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

func CrackHS256(tokenString string, worldlistPath string) (string, error) {
	file, err := os.Open(worldlistPath)
	if err != nil {
		return "", fmt.Errorf("error opening wordlist file: %w", err)
	}
	defer file.Close()

	lineCount := countLines(worldlistPath)
	if lineCount <= 0 {
		lineCount = 1
	}
	bar := progressbar.Default(int64(lineCount))

	numWorkers := runtime.NumCPU()
	if numWorkers < 1 {
		numWorkers = 1
	}

	secrets := make(chan string, numWorkers*2)
	result := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var tried int64

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case secret, ok := <-secrets:
					if !ok {
						return
					}
					atomic.AddInt64(&tried, 1)
					bar.Add(1)

					valid, err := VerifyHS256(tokenString, []byte(secret))
					if err != nil {
						// ignore errors for individual attempts
						continue
					}
					if valid {
						select {
						case result <- secret:
						default:
						}
						cancel()
						return
					}
				}
			}
		}()
	}

	go func() {
		defer func() {
			close(secrets)
		}()
		_, _ = file.Seek(0, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			secret := scanner.Text()
			select {
			case <-ctx.Done():
				return
			case secrets <- secret:
			}
		}
	}()

	var foundSecret string
	select {
	case s := <-result:
		foundSecret = s
		wg.Wait()
	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()
		return done
	}():
	}

	time.Sleep(50 * time.Millisecond)
	bar.Finish()

	if foundSecret != "" {
		fmt.Println("\nToken is valid with secret:", foundSecret)
		return foundSecret, nil
	}

	return "", fmt.Errorf("no valid secret found (tried %d words)", tried)
}

func countLines(filePath string) int {
	f, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
