package main

func main() {
	message := make(chan string)

	// Горутина (main) пытается получить сообщение из канала
	<-message // fatal error: all goroutines are asleep - deadlock!
}
