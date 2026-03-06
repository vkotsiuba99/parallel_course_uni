package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

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
