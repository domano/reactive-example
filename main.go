package main

import (
	"fmt"
	"github.com/briandowns/formatifier"
	"time"
)

func main() {
	// Observable / Output-Channel
	emitter := time.NewTicker(time.Second)

	// Observer: Creates a Log-Entry & Returns an Observable / Output-Channel
	logChan := logOberserver(emitter)

	// Oberserver: Transforms string to leet speak & Returns an Obserable / Output-Channel
	leetChan := morseOberserver(logChan)

	// Observer: Print Line
	printOberserver(leetChan)

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

func morseOberserver(input <-chan string) <-chan string {
	morseChan := make(chan string)
	go func() {
		for {
			select {
			case str, open := <-input:
				if !open {
					break
				}
				morse, err := formatifier.ToMorseCode(str)
				if err != nil {
					continue
				}
				morseChan <- morse
			}
		}
	}()
	return morseChan
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
