package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Heartbeats — это способ параллельных процессов сигнализировать
// о жизни внешним сторонам. Свое название они получили из анатомии
// человека, где сердцебиение означает для наблюдателя жизнь.
// Heartbeats существовали еще до Go и остаются в нем полезными.
//
// Существует два различных типа heartbeats:
//
// Heartbeats, возникающие в определенный интервал времени.
// Heartbeats, возникающие в начале единицы работы

type typedef int

const (
	doWork1 typedef = iota
	doWork2
	doWork3
	doWork4
	doWork5
)

var (
	activeDoWork typedef = doWork1
)

func main() {
	done := make(chan any)

	time.AfterFunc(10*time.Second, func() {
		close(done)
	})

	const timeout = 2 * time.Second

	switch activeDoWork {
	case doWork1:
		heartbeat, results := doWork1Fn(done, timeout/2)
		receive1(heartbeat, results, timeout)
	case doWork2:
		heartbeat, results := doWork2Fn(done)
		receive2(heartbeat, results)
	case doWork3:
		intSlice := []int{0, 1, 2, 3, 5}
		heartbeat, results := doWork3Fn(done, intSlice...)
		receive3(heartbeat, results, intSlice)
	case doWork4:
		intSlice := []int{0, 1, 2, 3, 5}
		heartbeat, results := doWork3Fn(done, intSlice...)
		<-heartbeat
		receive4(heartbeat, results, intSlice)
	case doWork5:
		intSlice := []int{0, 1, 2, 3, 5}
		const timeout = 2 * time.Second
		heartbeat, results := doWork4Fn(done, timeout/2, intSlice...)
		<-heartbeat
		receive5(heartbeat, results, timeout, intSlice)
	}
}

func doWork1Fn(done <-chan any, pulseInterval time.Duration) (<-chan any, <-chan time.Time) {
	heartbeat := make(chan any)
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)

		workGen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
				fmt.Println("without listener")
			}
		}
		_ = sendPulse

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()

	return heartbeat, results
}

func doWork2Fn(done <-chan any) (<-chan any, <-chan int) {
	heartbeatStream := make(chan any)
	workStream := make(chan int)

	go func() {
		defer close(heartbeatStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeatStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()

	return heartbeatStream, workStream
}

func doWork3Fn(done <-chan any, nums ...int) (<-chan any, <-chan int) {
	heartbeat := make(chan any)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()

	return heartbeat, intStream
}

func doWork4Fn(done <-chan any, pulseInterval time.Duration, nums ...int) (<-chan any, <-chan int) {
	heartbeat := make(chan any)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)
		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeat, intStream
}

func receive1(heartbeat <-chan any, results <-chan time.Time, timeout time.Duration) {
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results in second %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func receive2(heartbeat <-chan any, results <-chan int) {
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

func receive3(heartbeat <-chan any, results <-chan int, intSlice []int) {
	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				fmt.Printf("index %v: expected %v, but received %v\n", i, expected, r)
			}
		case <-time.After(1 * time.Second):
			fmt.Printf("test timed out")
		}
	}
}

func receive4(heartbeat <-chan any, results <-chan int, intSlice []int) {
	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			fmt.Printf("index %v: expected %vv, but received %v", i, expected, r)
		}
		i++
	}
}

func receive5(heartbeat <-chan any, results <-chan int, timeout time.Duration, intSlice []int) {
	i := 0
	for {
		select {
		case r, ok := <-results:
			if !ok {
				return
			} else if expected := intSlice[i]; r != expected {
				fmt.Printf("intex %v: expected %v, but received %v\n", i, expected, r)
			}
			i++
		case <-heartbeat:
		case <-time.After(timeout):
			fmt.Println("test timed out")
		}
	}
}
