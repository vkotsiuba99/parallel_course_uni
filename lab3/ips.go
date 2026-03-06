package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// IPC (Message Passing)
func runIPCMessagePassing() {
	ch := make(chan int64)
	var wg sync.WaitGroup
	wg.Add(2)

	// Main proc
	go func() {
		defer wg.Done()
		val := int64(rand.Intn(100))
		fmt.Printf("Generated number: %d\n", val)
		ch <- val
		time.Sleep(time.Second)
		response := <-ch
		fmt.Printf("Getting number back from support proc: %d\n", response)
	}()

	// Support proc
	go func() {
		defer wg.Done()
		received := <-ch
		fmt.Printf("Support proc number: %d\n", received)
		ch <- received
	}()
	wg.Wait()
}
