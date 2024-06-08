package main

import (
	"fmt"
	"sync"
)

func main() {
	myPool := &sync.Pool{
		New: func() any {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	// Если экземпляр пула не создан, тогда будет
	// вызвана функция New, определенная в пуле.
	myPool.Get()
	instance := myPool.Get()
	fmt.Println("instance", instance)

	// Здесь мы вставляем ранее полученный экземпляр обратно
	// в пул, это увеличивает количество экземпляров.
	myPool.Put(instance)

	// Когда этот вызов будет выполнен, мы будем повторно
	// использовать ранее выделенный экземпляр и поместим
	// его обратно в пул.
	myPool.Get()

	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() any {
			fmt.Println("new calc pool")

			numCalcsCreated++
			mem := make([]byte, 1024)

			return &mem
		},
	}

	fmt.Println("calPool.New", calcPool.New())

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	calcPool.Get()

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
