package main

import "sync"

type Account struct {
	ID      int64
	Balance int64
	Mu      sync.Mutex
}
