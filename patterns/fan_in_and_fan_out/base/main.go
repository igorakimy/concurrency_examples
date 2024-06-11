package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

type goFile struct {
	name string
	data string
}

func main() {
	si := []goFile{
		{name: "utils.go", data: "package utils\n\nfunc Util() {}"},
		{name: "helper.go", data: "package helper\n\nfunc Helper() {}"},
		{name: "misc.go", data: "package misc\n\nfunc Misc() {}"},
		{name: "various.go", data: "package various\n\nfunc Various() {}"},
	}

	input := make(chan goFile, len(si))
	res1 := make(chan string, len(si))
	res2 := make(chan string, len(si))
	res3 := make(chan string, len(si))

	chans := fanOut(input, 3)
	var g errgroup.Group

	g.Go(func() error {
		return scanSQLInjection(chans[0], res1)
	})

	g.Go(func() error {
		return scanTimingExploits(chans[1], res2)
	})

	g.Go(func() error {
		return scanAuth(chans[2], res3)
	})

	g.Go(func() error {
		for _, d := range si {
			input <- d
		}
		close(input)
		return nil
	})

	g.Go(func() error {
		res := fanIn(res1, res2, res3)
		for r := range res {
			fmt.Println(r)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("main: done")
}

func fanIn[T any](chans ...chan T) chan T {
	res := make(chan T)
	var g errgroup.Group

	for _, c := range chans {
		c := c

		g.Go(func() error {
			for s := range c {
				res <- s
			}
			return nil
		})
	}

	go func() {
		_ = g.Wait()
		close(res)
	}()

	return res
}

func fanOut[T any](ch chan T, n int) []chan T {
	chans := make([]chan T, 0, n)

	for i := 0; i < n; i++ {
		chans = append(chans, make(chan T, 1))
	}

	go func() {
		for item := range ch {
			for _, c := range chans {
				select {
				case c <- item:
				case <-time.After(100 * time.Millisecond):
				}
			}
		}

		for _, c := range chans {
			close(c)
		}
	}()

	return chans
}

func mockScan() string {
	if rand.Intn(100) > 90 {
		return "ALERT - vulnerability found"
	}
	return "OK - All Correct"
}

func scanSQLInjection(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("SQL injection scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}

func scanTimingExploits(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("Timing exploits scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}

func scanAuth(data <-chan goFile, res chan<- string) error {
	for d := range data {
		res <- fmt.Sprintf("Authentication scan: %s scanned, result: %s", d.name, mockScan())
	}
	close(res)
	return nil
}
