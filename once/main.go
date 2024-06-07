package main

import (
	"fmt"
	"sync"
)

// Once дает гарантию того, что нечто будет выполнено
// только один раз, даже среди нескольких горутин.

func main() {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}
