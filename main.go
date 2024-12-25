package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func fetchURL(url string, wg *sync.WaitGroup, results chan<- string, errors chan<- error) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		errors <- fmt.Errorf("error during receiving %s: %w", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errors <- fmt.Errorf("error during reading %s: %w", url, err)
		return
	}

	results <- fmt.Sprintf("URL: %s\nContent: %s\n", url, string(body))
}

func main() {
	file, err := os.Create("data.html")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	urls := []string{
		"https://metanit.com/go/tutorial/9.6.php",
		"https://golang.org",
		"https://openai.com",
	}

	var wg sync.WaitGroup
	results := make(chan string, len(urls))
	errors := make(chan error, len(urls))

	start := time.Now()

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg, results, errors)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	for {
		select {
		case result, ok := <-results:
			if ok {
				file.WriteString(result)
				fmt.Println("done")
			} else {
				results = nil
			}
		case err, ok := <-errors:
			if ok {
				log.Println("Error:", err)
			} else {
				errors = nil
			}
		}

		if results == nil && errors == nil {
			break
		}
	}

	fmt.Println("Ended in", time.Since(start))
}
