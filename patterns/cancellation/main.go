package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan any)
	strings := make(chan string)
	terminated := doWork(done, strings)

	go func() {
		for i := 1; i <= 3; i++ {
			strings <- fmt.Sprintf("%d", i)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	fmt.Println("initiate...")
	d := <-terminated
	fmt.Println(d)
	fmt.Println("Done.")
}

func doWork(done <-chan any, strings <-chan string) <-chan any {
	terminated := make(chan any)

	go func() {
		defer fmt.Println("doWork exited.")
		defer close(terminated)
		for {
			select {
			case s := <-strings:
				// Do something interesting
				fmt.Println(s)
			case <-done:
				return
			}
		}
	}()

	fmt.Println("doWork initiate...")

	return terminated
}
