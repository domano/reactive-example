package main

import (
	"fmt"
	"github.com/briandowns/formatifier"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// Observable / Output-Channel
	emitter := time.NewTicker(time.Second)

	// Observer: Creates a Log-Entry & Returns an Observable / Output-Channel
	logChan := logOberserver(emitter, 20)

	// Broadcast into 2 channels
	logChan1, logChan2 := broadcastObserver(logChan)


	// Oberserver: Transforms string to upper speak & Returns an Obserable / Output-Channel
	upperChan := upperOberserver(logChan1)

	// Oberserver: Transforms string to morse code & Returns an Obserable / Output-Channel
	morseChan := morseOberserver(logChan2)

	// Merge upper and morse channels
	mergedChan := mergeObserver(upperChan, morseChan)

	// Observer: Print Line
	printOberserver(mergedChan)

	// Block main channel to keep the program running
	select {}
}

func logOberserver(emitter *time.Ticker, numWorkers int) <-chan string {
	logChan := make(chan string)
	for i := 0; i < numWorkers; i++ {
		go log(emitter, logChan)
	}
	return logChan
}

func log(emitter *time.Ticker, logChan chan string) {
	for {
		timeEvent := <-emitter.C
		h, m, s := timeEvent.Clock()
		logChan <- fmt.Sprintf(
			"Hour: %d, Minute: %d, Second: %d",
			h, m, s)
		<-time.After(time.Duration(rand.Intn(6)) * time.Second)
	}
}

func broadcastObserver(in <-chan string) (out1, out2 chan string) {
	out1 = make(chan string)
	out2 = make(chan string)

	go broadcast(in, out1, out2)
	return out1, out2

}

func broadcast(in <-chan string, out1 chan<- string, out2 chan<- string){
	for {
		input := <-in
		out1 <- input
		out2 <- input
	}
}

func mergeObserver(in1, in2 <-chan string) (out chan string) {
	out = make(chan string)
	go merge(in1, in2, out)
	return out
}

func merge(in1, in2 <-chan string, out chan<- string, ) {
	for {
		select {
		case input := <-in1:
			out <- input
		case input := <-in2:
			out <- input
		}
	}
}

func morseOberserver(input <-chan string) <-chan string {
	morseChan := make(chan string)
	go morse(input, morseChan)
	return morseChan
}

func morse(input <-chan string, morseChan chan string) {
	for {
		str := <-input
		morse, err := formatifier.ToMorseCode(str)
		if err != nil {
			continue
		}
		morseChan <- morse
	}
}

func upperOberserver(input <-chan string) <-chan string {
	morseChan := make(chan string)
	go upper(input, morseChan)
	return morseChan
}

func upper(input <-chan string, upperChan chan string) {
	for {
		str := <-input
		upper := strings.ToUpper(str)
		upperChan <- upper
	}
}

func printOberserver(input <-chan string) {
	go func() {
		for {
			str := <-input
			println(str)
		}
	}()
}
