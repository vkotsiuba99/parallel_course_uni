package main

import (
	"fmt"
)

func main() {
	//targetDir := "./data"

	//var fileCount int64 = 20
	//_ = os.RemoveAll(targetDir)
	//
	//err := os.MkdirAll(targetDir, 0755)
	//if err != nil {
	//	fmt.Println("cannot create target dir")
	//	return
	//}
	//
	//fmt.Printf("fileCount: %d targetDir: %s...\n", fileCount, targetDir)
	//start := time.Now()
	//createRandomFiles(targetDir, fileCount)
	//fmt.Printf("execution time: %v\n", time.Since(start))

	fmt.Println("Lab1:")

	runCPUBound(100_000_000, false)
	runCPUBound(100_000_000, true)

	//runMemoryBound(10000, false)
	//runMemoryBound(10000, true)
	//
	//runIOBound(targetDir, false)
	//runIOBound(targetDir, true)
}
