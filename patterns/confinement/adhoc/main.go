package main

import "fmt"

func main() {
	data := []int{1, 2, 3, 4, 5}

	handleData := make(chan int)
	go loopData(handleData, data)

	for num := range handleData {
		fmt.Println(num)
	}
}

func loopData(handleData chan<- int, data []int) {
	defer close(handleData)
	for i := range data {
		handleData <- data[i]
	}
}
