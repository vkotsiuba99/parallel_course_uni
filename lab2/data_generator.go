package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func GenerateData(filename string, count int64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := int64(0); i < count; i++ {
		userID := rand.Int63n(1000)
		amount := rand.Float64()*500 + 10.0
		// UserID, amount, currency, date, type_of_product
		line := fmt.Sprintf("%d,%.2f,USD,2026-02-28,Electronics\n", userID, amount)
		writer.WriteString(line)
	}

	return writer.Flush()
}
