package main

import (
	"fmt"
	"math"
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

func factorize(n int64) []int64 {
	factors := []int64{}
	d := int64(2)
	temp := n

	for temp > 1 {
		if temp%d == 0 {
			factors = append(factors, d)
			temp /= d
			continue
		}

		d++
		if d*d > temp {
			if temp > 1 {
				factors = append(factors, temp)
				break
			}
		}
	}

	return factors
}

func runFactorization(numbers []int64, parallel bool) {
	start := time.Now()
	if !parallel {
		for _, n := range numbers {
			factorize(n)
		}
		fmt.Printf("Factorization (sequential): time: %v\n",
			time.Since(start))

		return
	}

	var wg sync.WaitGroup
	for _, n := range numbers {
		wg.Add(1)
		go func(num int64) {
			defer wg.Done()
			factorize(num)
		}(n)
	}
	wg.Wait()

	fmt.Printf("Factorization (parallel): time: %v\n",
		time.Since(start))
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}

	for i := int64(2); i <= int64(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func findPrimesInRange(start, end int64) []int64 {
	primes := []int64{}
	for i := start; i <= end; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}

	return primes
}

func runPrimesParallel(start, end int64, parallel bool) {
	startTime := time.Now()
	if !parallel {
		res := findPrimesInRange(start, end)
		fmt.Printf("Prime numbers (sequential): %d, time: %v\n",
			len(res), time.Since(startTime))
		return
	}

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
	chunkSize := (end - start) / int64(numCPU)
	results := make(chan int, numCPU)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := start + int64(i)*chunkSize
			e := s + chunkSize
			if i == numCPU-1 {
				e = end
			}
			res := findPrimesInRange(s, e)
			results <- len(res)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for count := range results {
		total += count
	}
	fmt.Printf("Prime number (parallel): %d, time: %v\n",
		total, time.Since(startTime))
}
