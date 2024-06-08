package main

func main() {
	message := make(chan string)

	// Горутина (main) пытается отправить сообщение в канал
	message <- "message" // fatal error: all goroutines are asleep - deadlock!
}
