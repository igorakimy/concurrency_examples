package main

import (
	"fmt"
	"sync"
)

// Mutex означает "взаимное исключение" и является
// способом защиты критических разделов вашей программы.

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}

func main() {
	counter := Counter{}
	counter.Increment()
	fmt.Println(counter.Value())
}
