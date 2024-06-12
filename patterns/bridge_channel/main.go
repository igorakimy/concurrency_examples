package main

import "fmt"

// С помощью этих шаблонов можно создать функцию,
// которая разбивает канал каналов на один канал.

func main() {
	done := make(chan any)
	defer close(done)

	for v := range bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}
}

func bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		// Этот цикл отвечает за извлечение каналов из chanStream
		// и передачу их во вложенный цикл для использования
		for {
			var stream <-chan any
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func genVals() <-chan <-chan any {
	chanStream := make(chan (<-chan any))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan any, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func orDone(done, c <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
