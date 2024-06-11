package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan any)
	myChan := make(chan any)

	go func() {
		done <- struct{}{}
	}()

	go func() {
		defer close(myChan)
		for i := 0; i < 1000; i++ {
			myChan <- i
			time.Sleep(2 * time.Second)
		}
	}()

	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}

func orDone(done, c <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				fmt.Println("finish")
				return
			case v, ok := <-c:
				if ok == false {
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
