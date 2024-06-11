package main

import "fmt"

// Fan-out — это термин, описывающий процесс запуска нескольких
// горутин для обработки входных данных конвейера, а fan-in — это термин,
// описывающий процесс объединения нескольких выходных данных в один канал.

type data int

func main() {
	work := []data{1, 2, 3, 4, 5}

	const numWorkers = 3

	wch := make(chan data, len(work))
	res := make(chan data, len(work))

	// fan-out, один входящий канал для всех горутин.
	for i := 0; i < numWorkers; i++ {
		go worker(wch, res)
	}

	for _, w := range work {
		fmt.Println("send to wch:", w)
		wch <- w
	}
	close(wch)

	// fan-in, один результирующий канал.
	for range work {
		w := <-res
		fmt.Println("receive from res:", w)
	}
}

func worker(wch <-chan data, res chan<- data) {
	for {
		w, ok := <-wch
		if !ok {
			return
		}

		w *= 2
		res <- w
	}
}
