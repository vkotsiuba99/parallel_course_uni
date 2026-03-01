package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Account struct {
	ID      int64
	Balance int64
	Mu      sync.Mutex
}

func GenerateAccountsFile(filename string, count int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := int64(0); i < count; i++ {
		balance := int64(rand.Intn(10000) + 1000)
		fmt.Fprintf(file, "%d,%d\n", i, balance)
	}
	return nil
}

func LoadAccounts(filename string) ([]*Account, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var accounts []*Account
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		id, _ := strconv.ParseInt(parts[0], 10, 64)
		bal, _ := strconv.ParseInt(parts[1], 10, 64)
		accounts = append(accounts, &Account{ID: id, Balance: bal})
	}
	return accounts, nil
}

// Race Condition
func transferRaceCondition(from, to *Account, amount int64, wg *sync.WaitGroup) {
	defer wg.Done()
	if from.Balance >= amount {
		tempFrom := from.Balance
		tempTo := to.Balance
		time.Sleep(time.Microsecond)
		from.Balance = tempFrom - amount
		to.Balance = tempTo + amount
	}
}

// Fixed (without Race Condition)
func transferFixed(from, to *Account, amount int64, wg *sync.WaitGroup) {
	defer wg.Done()

	if from.ID == to.ID {
		return
	}

	// sorted by ID for preventing Deadlock (Lock Ordering)
	var first, second *Account
	if from.ID < to.ID {
		first = from
		second = to
	} else {
		first = to
		second = from
	}

	first.Mu.Lock()
	second.Mu.Lock()
	defer second.Mu.Unlock()
	defer first.Mu.Unlock()

	if from.Balance >= amount {
		from.Balance -= amount
		to.Balance += amount
	}
}

func getTotalBalance(accounts []*Account) int64 {
	var total int64 = 0
	for _, acc := range accounts {
		total += acc.Balance
	}
	return total
}

// IPC (Message Passing)
func runIPCMessagePassing() {
	ch := make(chan int64)
	var wg sync.WaitGroup
	wg.Add(2)

	// Main proc
	go func() {
		defer wg.Done()
		val := int64(rand.Intn(100))
		fmt.Printf("Generated number: %d\n", val)
		ch <- val
		time.Sleep(time.Second)
		response := <-ch
		fmt.Printf("Getting number back from support proc: %d\n", response)
	}()

	// Support proc
	go func() {
		defer wg.Done()
		received := <-ch
		fmt.Printf("Support proc number: %d\n", received)
		ch <- received
	}()
	wg.Wait()
}

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
	fmt.Printf("Balance after після Condition: %d \n", getTotalBalance(accounts))

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
