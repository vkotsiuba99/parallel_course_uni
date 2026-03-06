package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	const accCount int64 = 100
	const transCount int64 = 1000
	filename := "accounts.txt"

	// Data generator
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		GenerateAccountsFile(filename, accCount)
	}

	accounts, _ := LoadAccounts(filename)
	initialTotal := getTotalBalance(accounts)
	fmt.Printf("Started balance: %d\n", initialTotal)

	// task 1.1 Race Condition
	var wg sync.WaitGroup
	fmt.Println("\nStart Race Condition (without sync)...")
	for i := int64(0); i < transCount; i++ {
		wg.Add(1)
		// random accounts
		f := accounts[rand.Int63n(accCount)]
		t := accounts[rand.Int63n(accCount)]
		go transferRaceCondition(f, t, 10, &wg)
	}
	wg.Wait()
	fmt.Printf("Balance after Condition: %d \n", getTotalBalance(accounts))

	// reload data
	accounts, _ = LoadAccounts(filename)

	// task 1.2: Fixed Version
	fmt.Println("\nStart Fixed Version (Mutex + ID Ordering)...")
	for i := int64(0); i < transCount; i++ {
		wg.Add(1)
		f := accounts[rand.Int63n(accCount)]
		t := accounts[rand.Int63n(accCount)]
		go transferFixed(f, t, 10, &wg)
	}
	wg.Wait()
	fmt.Printf("Balance after Fixed Version: %d \n\n\n", getTotalBalance(accounts))

	// task 2
	runIPCMessagePassing()
}

func getTotalBalance(accounts []*Account) int64 {
	var total int64 = 0
	for _, acc := range accounts {
		total += acc.Balance
	}
	return total
}
