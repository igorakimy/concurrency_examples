package main

import "fmt"

func main() {
	initStream := make(chan int)

	go func() {
		defer close(initStream)
		for i := 0; i <= 5; i++ {
			initStream <- i
		}
	}()

	// При перемещении по каналу for range прекращает цикл, если канал закрыт.
	for integer := range initStream {
		fmt.Printf("%v ", integer)
	}
}
