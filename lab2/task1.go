package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main2() {
	tagTypes := []string{"<div>", "<p>", "<a>"}
	docs := make([][]string, 10_000)
	for i := range docs {
		count := rand.Int63n(10_000)
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
	minP, maxP, sumP := statsForkJoinLegacy(arr)
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
