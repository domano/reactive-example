package main

import (
	"fmt"
	"time"
)

func main() {
	// Observable / Output-Channel
	emitter := time.NewTicker(time.Second)

	// Observer: Creates a Log-Entry & Returns an Observable / Output-Channel
	logChan := logOberserver(emitter)

	// Observer: Print Line
	printOberserver(logChan)

	// Block main channel to keep the program running
	select {}
}

func logOberserver(emitter *time.Ticker) <-chan string {
	logChan := make(chan string)
	go func() {
		for {
			select {
			case timeEvent, open := <-emitter.C:
				if !open {
					break
				}
				h, m, s := timeEvent.Clock()
				logChan <- fmt.Sprintf(
					"Hour: %d, Minute: %d, Second: %d",
					h, m, s)
			}
		}
	}()
	return logChan
}

func printOberserver(input <-chan string) {
	go func() {
		for {
			select {
			case str, open := <-input:
				if !open {
					break
				}
				println(str)
			}
		}
	}()
}
