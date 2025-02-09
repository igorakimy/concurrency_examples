package main

import (
	"fmt"
	"sync"
	"time"
)

// Starvation - это любая ситуация, когда параллельный процесс не может
// получить все ресурсы, необходимые для выполнения работы.

func main() {
	fmt.Println("vim-go")

	var wg sync.WaitGroup
	var sharedLock sync.Mutex
	const runtime = 1 * time.Second

	greedyWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Millisecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Millisecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Millisecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Millisecond)
			sharedLock.Unlock()

			count++
		}

		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}
