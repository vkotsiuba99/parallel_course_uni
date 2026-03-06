package main

import (
	"sync"
	"time"
)

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
