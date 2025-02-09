package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	done := make(chan any)
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}

		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func checkStatus(done <-chan any, urls ...string) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{Error: err, Response: resp}

			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()

	return results
}
