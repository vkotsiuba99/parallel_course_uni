package main

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
