package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// CPU-BOUND

// estimatePi Monte Carlo method
func estimatePi(iterations int64) float64 {
	var hits int64 = 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := int64(0); i < iterations; i++ {
		x, y := r.Float64(), r.Float64()
		if x*x+y*y <= 1 {
			hits++
		}
	}

	return 4 * float64(hits) / float64(iterations)
}

func runCPUBound(iterations int64, parallel bool) {
	start := time.Now()

	if !parallel {
		result := estimatePi(iterations)
		fmt.Printf("CPU-Bound (sequential): %f, time: %v\n",
			result, time.Since(start))

		return
	}

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
	results := make(chan float64, numCPU)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- estimatePi(iterations / int64(numCPU))
		}()
	}
	wg.Wait()
	close(results)

	var total float64
	for r := range results {
		total += r
	}
	fmt.Printf("CPU-Bound (parallel): %f, time: %v\n",
		total/float64(numCPU), time.Since(start))
}
