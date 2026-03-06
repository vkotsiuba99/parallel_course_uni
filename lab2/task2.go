package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseLine(line string) Transaction {
	parts := strings.Split(line, ",")
	uid, _ := strconv.ParseInt(parts[0], 10, 64)
	amt, _ := strconv.ParseFloat(parts[1], 64)
	return Transaction{UserID: uid, Amount: amt, Currency: parts[2]}
}

func main1() {
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
