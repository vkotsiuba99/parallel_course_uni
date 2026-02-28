package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// --- task 1.1: HTML tags (MAP-REDUCE) ---
func tagsSequential(docs [][]string) map[string]int64 {
	res := make(map[string]int64)
	for _, doc := range docs {
		for _, tag := range doc {
			res[tag]++
		}
	}
	return res
}

func tagsParallel(docs [][]string) map[string]int64 {
	results := make(chan map[string]int64, int64(len(docs)))
	var wg sync.WaitGroup

	for _, doc := range docs {
		wg.Add(1)
		go func(d []string) {
			defer wg.Done()
			local := make(map[string]int64)
			for _, tag := range d {
				local[tag]++
			}
			results <- local
		}(doc)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	final := make(map[string]int64)
	for res := range results {
		for k, v := range res {
			final[k] += v
		}
	}
	return final
}

// --- task 1.2: array stats (FORK-JOIN) ---
func statsSequential(data []int64) (int64, int64, int64) {
	minV, maxV, sumV := data[0], data[0], int64(0)
	for _, v := range data {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
		sumV += v
	}
	return minV, maxV, sumV
}

func statsForkJoin(data []int64) (int64, int64, int64) {
	if int64(len(data)) <= 50000 {
		return statsSequential(data)
	}

	mid := int64(len(data)) / 2
	resChan := make(chan struct{ min, max, sum int64 })

	// Fork
	go func() {
		minL, maxL, sumL := statsForkJoin(data[:mid])
		resChan <- struct{ min, max, sum int64 }{minL, maxL, sumL}
	}()

	minR, maxR, sumR := statsForkJoin(data[mid:])
	left := <-resChan // Join

	return min(left.min, minR), max(left.max, maxR), left.sum + sumR
}

// --- task 1.3: matrix (WORKER POOL) ---
func matrixSequential(size int64, A, B [][]int32) {
	for i := int64(0); i < size; i++ {
		for j := int64(0); j < size; j++ {
			var sum int32 = 0
			for k := int64(0); k < size; k++ {
				sum += A[i][k] * B[k][j]
			}
		}
	}
}

func matrixParallel(size int64, A, B [][]int32) {
	numWorkers := int64(4)
	var wg sync.WaitGroup
	rows := make(chan int64, size)

	for w := int64(0); w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range rows {
				for j := int64(0); j < size; j++ {
					var sum int32 = 0
					for k := int64(0); k < size; k++ {
						sum += A[r][k] * B[k][j]
					}
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

func main1() {
	tagTypes := []string{"<div>", "<p>", "<a>"}
	docs := make([][]string, 1000)
	for i := range docs {
		count := rand.Int63n(100)
		doc := make([]string, count)
		for j := int64(0); j < count; j++ {
			doc[j] = tagTypes[rand.Intn(len(tagTypes))]
		}
		docs[i] = doc
	}

	// 1
	fmt.Println("\ntask 1")
	startS := time.Now()
	resS := tagsSequential(docs)
	fmt.Printf("Sequential Tags: %v (Time: %v)\n", resS, time.Since(startS))

	startP := time.Now()
	resP := tagsParallel(docs)
	fmt.Printf("Parallel Tags: %v (Time: %v)\n", resP, time.Since(startP))

	time.Sleep(1 * time.Second)

	// 2
	fmt.Println("\ntask 2")
	arr := make([]int64, 1_000_000)
	for i := range arr {
		arr[i] = rand.Int63n(1000000)
	}

	startS = time.Now()
	minS, maxS, sumS := statsSequential(arr)
	fmt.Printf("Sequential Stats: Min:%d Max:%d Avg:%d (Time: %v)\n", minS, maxS, sumS/1000000, time.Since(startS))

	startP = time.Now()
	minP, maxP, sumP := statsForkJoin(arr)
	fmt.Printf("Parallel Fork-Join: Min:%d Max:%d Avg:%d (Time: %v)\n", minP, maxP, sumP/1000000, time.Since(startP))

	time.Sleep(1 * time.Second)

	// 3
	fmt.Println("\ntask 3")
	mSize := int64(500)
	A, B := make([][]int32, mSize), make([][]int32, mSize)
	for i := range A {
		A[i], B[i] = make([]int32, mSize), make([]int32, mSize)
	}

	startS = time.Now()
	matrixSequential(mSize, A, B)
	fmt.Printf("Sequential Matrix: Time: %v\n", time.Since(startS))

	startP = time.Now()
	matrixParallel(mSize, A, B)
	fmt.Printf("Parallel Matrix: Time: %v\n", time.Since(startP))

	time.Sleep(1 * time.Second)
}
