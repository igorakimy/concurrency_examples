package main

import (
	"fmt"
	"sync"
	"time"
)

// Было бы лучше, если бы у горутины был какой-то способ эффективно
// спать до тех пор, пока ей не будет подан сигнал о пробуждении и
// проверке ее состояния. Именно это и делает для нас тип Cond.
//
// Cond и Broadcast — это метод, который обеспечивает уведомление горутин,
// заблокированных при вызове Wait, о том, что условие сработало.

func main() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]any, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()

		queue = queue[1:]

		fmt.Println("Removed from queue")

		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()

		for len(queue) == 2 {
			c.Wait()
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})

		go removeFromQueue(1 * time.Second)

		c.L.Unlock()
	}
}
