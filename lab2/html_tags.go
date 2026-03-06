package main

import "sync"

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
