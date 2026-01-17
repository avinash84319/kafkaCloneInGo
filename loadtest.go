//go:build loadtest

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Simple payload matching the API contract
type Payload struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
	Topic string `json:"topic"`
}

// randomString returns an alphanumeric string of length n
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Seed RNG
	rand.Seed(time.Now().UnixNano())

	url := "http://localhost:8080/sendMessage"
	total := 100000

	// Optional overrides via env vars
	if v := os.Getenv("LOADTEST_TOTAL"); v != "" {
		var n int
		fmt.Sscanf(v, "%d", &n)
		if n > 0 {
			total = n
		}
	}
	if v := os.Getenv("LOADTEST_URL"); v != "" {
		url = v
	}

	// Optional: comma-separated topics via env var, else default
	var topics []string
	if v := os.Getenv("LOADTEST_TOPICS"); v != "" {
		topics = strings.Split(v, ",")
		for i := range topics {
			topics[i] = strings.TrimSpace(topics[i])
		}
	} else {
		topics = []string{"testing", "alpha", "beta", "gamma"}
	}

	concurrency := runtime.NumCPU() * 4 // moderate concurrency; still launches total goroutines
	client := &http.Client{Timeout: 10 * time.Second}

	var wg sync.WaitGroup
	wg.Add(total)

	var success int64
	var fail int64

	start := time.Now()

	// Launch total goroutines (1 request each) to fully parallelize
	for i := 0; i < total; i++ {
		go func(i int) {
			defer wg.Done()
			// Generate random payload
			randomKey := rand.Intn(1_000_000) + 1            // 1..1e6
			randomVal := randomString(rand.Intn(32-8+1) + 8) // 8..32 chars
			topic := topics[rand.Intn(len(topics))]

			payload := Payload{
				Key:   randomKey,
				Value: randomVal,
				Topic: topic,
			}
			body, _ := json.Marshal(payload)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				atomic.AddInt64(&fail, 1)
				fmt.Printf("FAIL (req build) id=%d err=%v\n", i, err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				atomic.AddInt64(&fail, 1)
				fmt.Printf("FAIL (http) id=%d err=%v\n", i, err)
				return
			}
			// Read body then close
			respBody, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				atomic.AddInt64(&success, 1)
			} else {
				atomic.AddInt64(&fail, 1)
				// Print status and a trimmed body to avoid huge logs
				trim := string(respBody)
				if len(trim) > 512 {
					trim = trim[:512] + "..."
				}
				fmt.Printf("FAIL (status) id=%d status=%d body=%q\n", i, resp.StatusCode, trim)
			}
		}(i)
	}

	// Basic worker pool knob to avoid too much scheduling contention
	_ = concurrency // left in case future tuning is needed

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("Load test summary:")
	fmt.Printf("Total: %d, Success: %d, Fail: %d\n", total, success, fail)
	fmt.Printf("Elapsed: %v\n", elapsed)
	if elapsed > 0 {
		qps := float64(success) / elapsed.Seconds()
		fmt.Printf("Approx QPS: %.2f\n", qps)
		avgLatencyMs := (elapsed.Seconds() * 1000.0) / float64(total)
		fmt.Printf("Avg latency per request (amortized): %.2f ms\n", avgLatencyMs)
	}
}
