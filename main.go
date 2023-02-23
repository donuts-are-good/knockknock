package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func loadTest(url string, duration time.Duration, wg *sync.WaitGroup) {
	startTime := time.Now()
	endTime := startTime.Add(duration)

	client := &http.Client{}

	var successCount, failureCount, totalTime, highestTime, lowestTime int64
	highestTime = -1

	for time.Now().Before(endTime) {
		req, _ := http.NewRequest("GET", url, nil)
		startTime := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			failureCount++
		} else {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				successCount++
				responseTime := time.Since(startTime).Milliseconds()
				totalTime += responseTime
				if responseTime > highestTime {
					highestTime = responseTime
				}
				if lowestTime == 0 || responseTime < lowestTime {
					lowestTime = responseTime
				}
			} else {
				failureCount++
			}
		}
	}

	averageTime := totalTime / successCount

	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Requests: %v\n", successCount+failureCount)
	fmt.Printf("Successes: %v\n", successCount)
	fmt.Printf("Failures: %v\n", failureCount)
	fmt.Printf("Average response time: %v ms\n", averageTime)
	fmt.Printf("Highest response time: %v ms\n", highestTime)
	fmt.Printf("Lowest response time: %v ms\n", lowestTime)

	wg.Done()
}

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	}))

	var wg sync.WaitGroup

	wg.Add(1)
	go loadTest(server.URL, 60*time.Second, &wg)

	wg.Wait()
	server.Close()
}
