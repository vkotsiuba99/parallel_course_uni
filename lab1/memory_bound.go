package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MEMORY-BOUND (Transpose of a matrix)

// runMemoryBound create matrix 10000x10000 and transpose
// limited by speed of access of RAM
func runMemoryBound(size int, parallel bool) {
	matrix := make([][]int64, size)
	for i := range matrix {
		matrix[i] = make([]int64, size)
	}
	start := time.Now()

	if !parallel {
		for i := 0; i < size; i++ {
			for j := i + 1; j < size; j++ {
				matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
			}
		}

		fmt.Printf("Memory-Bound (sequential): time: %v\n", time.Since(start))
		return
	}

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	chunkSize := size / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(sRow, eRow int) {
			defer wg.Done()
			for i := sRow; i < eRow; i++ {
				for j := i + 1; j < size; j++ {
					matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
				}
			}
		}(w*chunkSize, (w+1)*chunkSize)
	}
	wg.Wait()

	fmt.Printf("Memory-Bound (parallel): time: %v\n", time.Since(start))
}
