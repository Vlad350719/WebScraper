package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestFetchURL(t *testing.T) {
    var wg sync.WaitGroup
    results := make(chan string, 1)
    errors := make(chan error, 1)

    // Create a mock
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, client")
    }))
    defer server.Close()

    wg.Add(1)
    go fetchURL(server.URL, &wg, results, errors)

    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    // Check results
    select {
    case result := <-results:
        if !strings.Contains(result, "Hello, client") {
            t.Errorf("Expected 'Hello, client' in result, got %q", result)
        }
    case err := <-errors:
        t.Errorf("Unexpected error: %v", err)
    }
}

func setupMockServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, client")
    }))
}

func BenchmarkFetchURL(b *testing.B) {
    // Setup mock
    server := setupMockServer()
    defer server.Close()

    var wg sync.WaitGroup
    results := make(chan string, 1)
    errors := make(chan error, 1)

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        wg.Add(1)
        go fetchURL(server.URL, &wg, results, errors)

        go func() {
            wg.Wait()
            close(results)
            close(errors)
        }()

        for {
            select {
            case <-results:
                // Handle result
            case <-errors:
                // Handle error
            }

            if len(results) == 0 && len(errors) == 0 {
                break
            }
        }
    }

    b.StopTimer()
}
