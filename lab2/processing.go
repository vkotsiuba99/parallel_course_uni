package main

import (
	"bufio"
	"os"
	"sync"
)

func SequentialProcess(filename string) float64 {
	file, _ := os.Open(filename)
	defer file.Close()

	var total float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tx := parseLine(scanner.Text())
		// converting and cashback
		val := tx.Amount * 40.5
		if tx.UserID > 500 {
			val *= 0.8
		}
		total += val
	}
	return total
}

func ProducerConsumerProcess(filename string) float64 {
	jobs := make(chan Transaction, 100)
	results := make(chan float64, 100)
	var wg sync.WaitGroup

	// Producer
	go func() {
		file, _ := os.Open(filename)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jobs <- parseLine(scanner.Text())
		}
		close(jobs)
	}()

	// Consumers
	numWorkers := int32(4)
	for i := int32(0); i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tx := range jobs {
				val := tx.Amount * 40.5
				if tx.UserID > 500 {
					val *= 0.8
				}
				results <- val
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var total float64
	for r := range results {
		total += r
	}
	return total
}

func PipelineProcess(filename string) float64 {
	// Stage 1: Reading
	source := func() <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			file, _ := os.Open(filename)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				out <- parseLine(scanner.Text())
			}
			close(out)
			file.Close()
		}()
		return out
	}

	// Stage 2: Converter (USD -> UAH)
	convert := func(in <-chan Transaction) <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			for tx := range in {
				tx.Amount *= 40.5
				out <- tx
			}
			close(out)
		}()
		return out
	}

	// Stage 3: Cashback
	cashback := func(in <-chan Transaction) <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			for tx := range in {
				if tx.UserID > 500 {
					tx.Amount *= 0.8
				}
				out <- tx
			}
			close(out)
		}()
		return out
	}

	finalPipe := cashback(convert(source()))

	var total float64
	for tx := range finalPipe {
		total += tx.Amount
	}

	return total
}
