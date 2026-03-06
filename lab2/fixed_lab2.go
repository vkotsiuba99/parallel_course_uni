package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- Data Structures ---

type StatsResult struct {
	Min, Max, Sum int64
}

// ============================================================================
// TASK 1: Pattern Research (Tags, Statistics, Matrix)
// ============================================================================

// 1.1 Tag Counting using Map-Reduce
func tagsMapReduce(docs [][]string, chunks int64) map[string]int64 {
	chunkSize := int64(len(docs)) / chunks
	results := make(chan map[string]int64, chunks)
	var wg sync.WaitGroup

	for i := int64(0); i < chunks; i++ {
		wg.Add(1)
		go func(idx int64) {
			defer wg.Done()
			local := make(map[string]int64)
			start := idx * chunkSize
			end := (idx + 1) * chunkSize
			if idx == chunks-1 {
				end = int64(len(docs))
			}

			for _, doc := range docs[start:end] {
				for _, tag := range doc {
					local[tag]++
				}
			}
			results <- local
		}(i)
	}
	wg.Wait()
	close(results)

	final := make(map[string]int64)
	for res := range results {
		for k, v := range res {
			final[k] += v
		}
	}
	return final
}

// 1.2 Array Statistics using Fork-Join
func statsForkJoin(arr []int64, threshold int64) StatsResult {
	if int64(len(arr)) <= threshold {
		res := StatsResult{Min: arr[0], Max: arr[0], Sum: 0}
		for _, v := range arr {
			if v < res.Min {
				res.Min = v
			}
			if v > res.Max {
				res.Max = v
			}
			res.Sum += v
		}
		return res
	}

	mid := int64(len(arr)) / 2
	leftChan := make(chan StatsResult)

	go func() { leftChan <- statsForkJoin(arr[:mid], threshold) }()
	right := statsForkJoin(arr[mid:], threshold)
	left := <-leftChan

	return StatsResult{
		Min: min(left.Min, right.Min),
		Max: max(left.Max, right.Max),
		Sum: left.Sum + right.Sum,
	}
}

// 1.3 Matrix Multiplication using Worker Pool
func matrixWorkerPool(size int64, workers int64) {
	// Mock matrices (A and B)
	rows := make(chan int64, size)
	var wg sync.WaitGroup

	for w := int64(0); w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range rows {
				// Mock multiplication logic for row R
				for j := int64(0); j < size; j++ {
					var sum int32 = 0
					for k := int64(0); k < size; k++ {
						sum += int32(r+k) * int32(k+j) // Dummy operation
					}
					_ = sum
				}
			}
		}()
	}
	for i := int64(0); i < size; i++ {
		rows <- i
	}
	close(rows)
	wg.Wait()
}

// ============================================================================
// TASK 2: Financial Transactions (File-based)
// ============================================================================

func generateTxFile(filename string, count int64) {
	file, _ := os.Create(filename)
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i := int64(0); i < count; i++ {
		// UserID, Amount, Currency, Date, Type
		fmt.Fprintf(writer, "%d,%.2f,USD,2026-03-06,Item\n", rand.Int63n(1000), rand.Float64()*500)
	}
	writer.Flush()
}

func parseTx(line string) Transaction {
	p := strings.Split(line, ",")
	uid, _ := strconv.ParseInt(p[0], 10, 64)
	amt, _ := strconv.ParseFloat(p[1], 64)
	return Transaction{UserID: uid, Amount: amt, Currency: p[2]}
}

func runPipeline(filename string) float64 {
	source := func() <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			f, _ := os.Open(filename)
			s := bufio.NewScanner(f)
			for s.Scan() {
				out <- parseTx(s.Text())
			}
			close(out)
			f.Close()
		}()
		return out
	}

	convert := func(in <-chan Transaction) <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			for tx := range in {
				tx.Amount *= 40.0
				out <- tx
			} // USD to UAH
			close(out)
		}()
		return out
	}

	cashback := func(in <-chan Transaction) <-chan Transaction {
		out := make(chan Transaction)
		go func() {
			for tx := range in {
				if tx.UserID > 500 {
					tx.Amount *= 0.8
				} // 20% cashback
				out <- tx
			}
			close(out)
		}()
		return out
	}

	var total float64
	for tx := range cashback(convert(source())) {
		total += tx.Amount
	}
	return total
}

// ============================================================================
// RESEARCH & STATISTICS
// ============================================================================

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("=== Parallel Programming Research (2026) ===")

	// Data Prep
	arrSize := int64(5_000_000)
	arr := make([]int64, arrSize)
	for i := range arr {
		arr[i] = rand.Int63n(1000)
	}

	// 1. Benchmark Sequential Baseline
	startSeq := time.Now()
	statsForkJoin(arr, arrSize+1) // Threshold > size forces sequential
	seqTime := time.Since(startSeq)
	fmt.Printf("Sequential Time (Baseline): %v\n\n", seqTime)

	// 2. Detailed Performance Table
	threadCounts := []int64{1, 2, 4, 8, 16}
	fmt.Printf("%-15s | %-8s | %-12s | %-10s | %-10s\n", "Pattern", "Threads", "Time", "Speedup", "Efficiency")
	fmt.Println(strings.Repeat("-", 65))

	for _, tc := range threadCounts {
		start := time.Now()
		statsForkJoin(arr, arrSize/tc) // Dynamic threshold for fork-join
		dur := time.Since(start)

		speedup := float64(seqTime.Nanoseconds()) / float64(dur.Nanoseconds())
		efficiency := speedup / float64(tc)

		fmt.Printf("%-15s | %-8d | %-12v | %-10.2fx | %-10.2f\n",
			"Fork-Join", tc, dur.Truncate(time.Microsecond), speedup, efficiency)
	}

	// Task 2 Execution
	fmt.Println("\n=== Processing Financial Transactions (Pipeline) ===")
	generateTxFile("tx.txt", 10000)
	total := runPipeline("tx.txt")
	fmt.Printf("Total Sum After Processing: %.2f UAH\n", total)
}
