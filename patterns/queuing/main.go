package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Буферизованный канал — это тип очереди. Преждевременное добавление
// очереди может скрыть проблемы синхронизации, такие как взаимоблокировки;
// мы можем использовать очередь, чтобы установить ограничение на обработку,
// в этом процессе, когда limit <- struct{}{} заполнен, очередь ожидает
// освобождения <-limit, если мы их удалим, одновременно будет создано 50 горутин.

func main() {
	var wg sync.WaitGroup
	limit := make(chan any, runtime.NumCPU())

	fmt.Printf("Started, Limit %d\n", cap(limit))

	workers(limit, &wg)
	wg.Wait()

	fmt.Println("Finished")
}

func workers(l chan any, wg *sync.WaitGroup) {
	for i := 0; i <= 50; i++ {
		i := i

		l <- struct{}{}
		wg.Add(1)

		go func(x int, w *sync.WaitGroup) {
			defer wg.Done()

			time.Sleep(1 * time.Second)
			fmt.Printf("Process %d\n", x)

			<-l
		}(i, wg)
	}
}
