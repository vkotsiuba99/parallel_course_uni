package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Transaction struct {
	UserID   int64
	Amount   float64
	Currency string
}

func GenerateData(filename string, count int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := int64(0); i < count; i++ {
		userID := rand.Int63n(1000)
		amount := rand.Float64()*500 + 10.0
		// UserID, amount, currency, date, type_of_product
		line := fmt.Sprintf("%d,%.2f,USD,2026-02-28,Electronics\n", userID, amount)
		writer.WriteString(line)
	}

	return writer.Flush()
}

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

func parseLine(line string) Transaction {
	parts := strings.Split(line, ",")
	uid, _ := strconv.ParseInt(parts[0], 10, 64)
	amt, _ := strconv.ParseFloat(parts[1], 64)
	return Transaction{UserID: uid, Amount: amt, Currency: parts[2]}
}

func main() {
	fileName := "transactions.txt"
	//count := int64(10000)
	//
	//fmt.Printf("Generation %d txs in file ...\n", count)
	//err := GenerateData(fileName, count)
	//if err != nil {
	//	log.Fatalf("Error generating data: %s", err)
	//}

	start := time.Now()
	resSeq := SequentialProcess(fileName)
	fmt.Printf("Sequential: %.2f (Time: %v)\n", resSeq, time.Since(start))

	start = time.Now()
	resPC := ProducerConsumerProcess(fileName)
	fmt.Printf("Producer-Consumer: %.2f (Time: %v)\n", resPC, time.Since(start))

	start = time.Now()
	resPipe := PipelineProcess(fileName)
	fmt.Printf("Pipeline: %.2f (Time: %v)\n", resPipe, time.Since(start))
}
