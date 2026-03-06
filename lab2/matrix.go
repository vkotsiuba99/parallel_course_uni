package main

import "sync"

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
